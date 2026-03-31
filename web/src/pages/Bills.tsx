import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, message, Popconfirm, Select, DatePicker, InputNumber, Tag, Pagination } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, PrinterOutlined, HistoryOutlined, SearchOutlined, ReloadOutlined, EyeOutlined } from '@ant-design/icons';
import type { Bill, Lease } from '../types/api';
import { billApi, type BillQueryParams, type BillsQueryResult } from '../api/bill';
import { leaseApi } from '../api/lease';
import dayjs from 'dayjs';
import OperationLogModal from '../components/OperationLogModal';

const { Option } = Select;

const typeMap: Record<string, string> = {
  rent: '租金',
  checkout: '退租结算',
};

const typeColorMap: Record<string, string> = {
  rent: 'blue',
  checkout: 'warning',
};

const statusMap: Record<string, string> = {
  pending: '待到账',
  paid: '已到账',
};

const statusColorMap: Record<string, string> = {
  pending: 'warning',
  paid: 'success',
};

const Bills: React.FC = () => {
  const [bills, setBills] = useState<Bill[]>([]);
  const [leases, setLeases] = useState<Lease[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [editingBill, setEditingBill] = useState<Bill | null>(null);
  const [viewingBill, setViewingBill] = useState<Bill | null>(null);
  const [form] = Form.useForm();
  const [operationLogVisible, setOperationLogVisible] = useState(false);
  const [currentBill, setCurrentBill] = useState<Bill | null>(null);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [queryForm] = Form.useForm();
  const [selectedType, setSelectedType] = useState<string>('');

  const fetchBills = async (params?: BillQueryParams) => {
    setLoading(true);
    try {
      const queryParams = {
        ...params,
        offset: (page - 1) * pageSize,
        limit: pageSize,
      };
      const data: BillsQueryResult = await billApi.list(queryParams);
      setBills(data.items);
      setTotal(data.total);
    } catch (error) {
      message.error('获取账单列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchLeases = async () => {
    try {
      const data = await leaseApi.list();
      setLeases(data.items || []);
    } catch (error) {
      message.error('获取租约列表失败');
    }
  };

  useEffect(() => {
    fetchBills();
    fetchLeases();
  }, [page, pageSize]);

  const handleQuery = async () => {
    const values = await queryForm.validateFields();
    setPage(1);
    fetchBills(values);
  };

  const handleReset = () => {
    queryForm.resetFields();
    setPage(1);
    fetchBills();
  };

  const handlePageChange = (pageNum: number, pageSizeNum: number) => {
    setPage(pageNum);
    setPageSize(pageSizeNum);
  };

  const handleCreate = () => {
    setEditingBill(null);
    setSelectedType('');
    form.resetFields();
    setModalVisible(true);
  };

  const handleTypeChange = (value: string) => {
    setSelectedType(value);
  };

  const handleEdit = (bill: Bill) => {
    setEditingBill(bill);
    form.setFieldsValue({
      ...bill,
      paidAt: bill.paidAt ? dayjs(bill.paidAt) : null,
      refundDepositAmount: bill.refundDepositAmount || 0,
    });
    setModalVisible(true);
  };

  const handlePrint = async (bill: Bill) => {
    try {
      await billApi.printReceipt(bill.id);
      message.success('打印成功');
    } catch (error) {
      message.error('打印失败');
    }
  };

  const handleConfirmArrival = async (bill: Bill) => {
    try {
      await billApi.confirmArrival(bill.id);
      message.success('到账确认成功');
      fetchBills();
    } catch (error) {
      message.error('到账确认失败');
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await billApi.delete(id);
      message.success('删除成功');
      fetchBills();
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleViewDetail = (bill: Bill) => {
    setViewingBill(bill);
    setDetailModalVisible(true);
  };

  const handleViewOperationLogs = (bill: Bill) => {
    setCurrentBill(bill);
    setOperationLogVisible(true);
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      const formattedValues = {
        ...values,
        paidAt: values.paidAt ? values.paidAt.format('YYYY-MM-DD') : null,
      };
      if (editingBill) {
        await billApi.update(editingBill.id, formattedValues);
        message.success('更新成功');
      } else {
        await billApi.create(formattedValues);
        message.success('创建成功');
      }
      setModalVisible(false);
      fetchBills();
    } catch (error) {
      message.error('操作失败');
    }
  };

  const formatAmount = (amount: number) => {
    return `¥${(amount / 100).toFixed(2)}`;
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    {
      title: '租约',
      dataIndex: 'leaseId',
      key: 'leaseId',
      render: (leaseId: string) => leases.find(l => l.id === leaseId)?.id || leaseId,
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => (
        <Tag color={typeColorMap[type] || 'default'}>
          {typeMap[type] || type}
        </Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={statusColorMap[status] || 'default'}>
          {statusMap[status] || status}
        </Tag>
      ),
    },
    {
      title: '金额明细',
      key: 'amountDetail',
      width: 350,
      render: (_: any, record: Bill) => {
        if (record.type === 'checkout') {
          // 退租结算账单 - 显示明细（所有项目都显示，即使是0）
          const refundRent = Math.abs(record.rentAmount || 0);
          const refundDeposit = record.refundDepositAmount || 0;
          const water = record.waterAmount || 0;
          const electric = record.electricAmount || 0;
          const other = record.otherAmount || 0;

          return (
            <div style={{ fontSize: '12px' }}>
              <div style={{ color: '#1890ff' }}>退还租金: {formatAmount(refundRent)}</div>
              <div style={{ color: '#1890ff' }}>退还押金: {formatAmount(refundDeposit)}</div>
              <div style={{ color: '#faad14' }}>水费: {formatAmount(water)}</div>
              <div style={{ color: '#faad14' }}>电费: {formatAmount(electric)}</div>
              <div style={{ color: '#faad14' }}>其他: {formatAmount(other)}</div>
              <div style={{ fontWeight: 'bold', marginTop: '4px' }}>
                合计: {formatAmount(record.amount)}
              </div>
            </div>
          );
        }
        // 租金账单 - 也显示详细明细
        const rent = record.rentAmount || 0;
        const water = record.waterAmount || 0;
        const electric = record.electricAmount || 0;
        const other = record.otherAmount || 0;

        // 检查是否有明细金额
        const hasDetails = rent > 0 || water > 0 || electric > 0 || other > 0;

        if (hasDetails) {
          return (
            <div style={{ fontSize: '12px' }}>
              {rent > 0 && <div style={{ color: '#1890ff' }}>租金: {formatAmount(rent)}</div>}
              {water > 0 && <div style={{ color: '#faad14' }}>水费: {formatAmount(water)}</div>}
              {electric > 0 && <div style={{ color: '#faad14' }}>电费: {formatAmount(electric)}</div>}
              {other > 0 && <div style={{ color: '#faad14' }}>其他: {formatAmount(other)}</div>}
              <div style={{ fontWeight: 'bold', marginTop: '4px' }}>
                合计: {formatAmount(record.amount)}
              </div>
            </div>
          );
        }
        // 没有明细的普通账单，只显示总金额
        return formatAmount(record.amount);
      },
    },
    { title: '到账时间', dataIndex: 'paidAt', key: 'paidAt' },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: Bill) => (
        <div>
          <Button
            icon={<EyeOutlined />}
            onClick={() => handleViewDetail(record)}
            style={{ marginRight: 8 }}
            size="small"
          >
            详情
          </Button>
          <Button
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
            style={{ marginRight: 8 }}
            size="small"
          >
            编辑
          </Button>
          <Button
            icon={<PrinterOutlined />}
            onClick={() => handlePrint(record)}
            style={{ marginRight: 8 }}
            size="small"
          >
            打印
          </Button>
          {record.status === 'pending' && (
            <Popconfirm
              title="确认到账?"
              onConfirm={() => handleConfirmArrival(record)}
              okText="确定"
              cancelText="取消"
            >
              <Button
                type="primary"
                size="small"
                style={{ marginRight: 8 }}
              >
                到账确认
              </Button>
            </Popconfirm>
          )}
          <Button
            icon={<HistoryOutlined />}
            onClick={() => handleViewOperationLogs(record)}
            style={{ marginRight: 8 }}
            size="small"
          >
            操作日志
          </Button>
          <Popconfirm
            title="确定删除?"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button icon={<DeleteOutlined />} danger size="small">
              删除
            </Button>
          </Popconfirm>
        </div>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h1>账单管理</h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新增账单
        </Button>
      </div>

      {/* 查询表单 */}
      <Form form={queryForm} layout="inline" style={{ marginBottom: 16 }}>
        <Form.Item name="type" label="类型">
          <Select placeholder="请选择类型" style={{ width: 120 }}>
            <Option value="rent">租金</Option>
            <Option value="checkout">退租结算</Option>
          </Select>
        </Form.Item>
        <Form.Item name="status" label="状态">
          <Select placeholder="请选择状态" style={{ width: 120 }}>
            <Option value="pending">待到账</Option>
            <Option value="paid">已到账</Option>
          </Select>
        </Form.Item>
        <Form.Item name="month" label="月份">
          <DatePicker
            picker="month"
            placeholder="请选择月份"
            style={{ width: 150 }}
          />
        </Form.Item>
        <Form.Item>
          <Button type="primary" icon={<SearchOutlined />} onClick={handleQuery} loading={loading}>
            查询
          </Button>
        </Form.Item>
        <Form.Item>
          <Button icon={<ReloadOutlined />} onClick={handleReset} loading={loading}>
            重置
          </Button>
        </Form.Item>
      </Form>

      <Table
        columns={columns}
        dataSource={bills}
        rowKey="id"
        loading={loading}
        pagination={false}
        scroll={{ x: 1800 }}
      />

      <div style={{ marginTop: 16, display: 'flex', justifyContent: 'flex-end' }}>
        <Pagination
          current={page}
          pageSize={pageSize}
          total={total}
          onChange={handlePageChange}
          showSizeChanger
          showQuickJumper
          showTotal={(total) => `共 ${total} 条记录`}
        />
      </div>

      <Modal
        title={editingBill ? '编辑账单' : '新增账单'}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={() => setModalVisible(false)}
        width={700}
      >
        <Form form={form} layout="vertical">
          {!editingBill && (
            <Form.Item
              name="leaseId"
              label="租约"
              rules={[{ required: true, message: '请选择租约' }]}
            >
              <Select placeholder="请选择租约">
                {leases.map(lease => (
                  <Option key={lease.id} value={lease.id}>
                    {lease.id}
                  </Option>
                ))}
              </Select>
            </Form.Item>
          )}
          {editingBill ? (
            <Form.Item
              name="type"
              label="类型"
            >
              <Select disabled>
                <Option value="rent">租金</Option>
                <Option value="checkout">退租结算</Option>
              </Select>
            </Form.Item>
          ) : (
            <Form.Item
              name="type"
              label="类型"
              rules={[{ required: true, message: '请选择类型' }]}
              initialValue="rent"
            >
              <Select placeholder="请选择类型" onChange={handleTypeChange}>
                <Option value="rent">租金</Option>
              </Select>
            </Form.Item>
          )}

          {/* 租金账单 - 显示所有金额字段 */}
          {(!editingBill || editingBill?.type === 'rent') && (
            <>
              <Form.Item
                name="rentAmount"
                label="租金（分）"
                rules={[{ required: true, message: '请输入租金' }]}
              >
                <InputNumber style={{ width: '100%' }} placeholder="请输入租金（分）" />
              </Form.Item>
              <Form.Item
                name="waterAmount"
                label="水费（分）"
              >
                <InputNumber style={{ width: '100%' }} placeholder="请输入水费（分）" />
              </Form.Item>
              <Form.Item
                name="electricAmount"
                label="电费（分）"
              >
                <InputNumber style={{ width: '100%' }} placeholder="请输入电费（分）" />
              </Form.Item>
              <Form.Item
                name="otherAmount"
                label="其他金额（分）"
              >
                <InputNumber style={{ width: '100%' }} placeholder="请输入其他金额（分）" />
              </Form.Item>
            </>
          )}

          {/* 退租结算账单 - 显示所有金额字段 */}
          {editingBill?.type === 'checkout' && (
            <>
              <div style={{ marginBottom: 16, padding: 16, border: '1px solid #d9d9d9', borderRadius: 8 }}>
                <h4 style={{ margin: '0 0 16px 0', color: '#1890ff' }}>退还金额</h4>
                <Form.Item
                  name="rentAmount"
                  label="退还租金（分，负数表示退还）"
                >
                  <InputNumber style={{ width: '100%' }} placeholder="请输入退还租金（分，负数表示退还）" />
                </Form.Item>
                <Form.Item
                  name="refundDepositAmount"
                  label="退还押金（分）"
                >
                  <InputNumber style={{ width: '100%' }} placeholder="请输入退还押金（分）" min={0} />
                </Form.Item>
              </div>
              <div style={{ marginBottom: 16, padding: 16, border: '1px solid #d9d9d9', borderRadius: 8 }}>
                <h4 style={{ margin: '0 0 16px 0', color: '#52c41a' }}>收取费用</h4>
                <Form.Item
                  name="waterAmount"
                  label="水费（分）"
                >
                  <InputNumber style={{ width: '100%' }} placeholder="请输入水费（分）" min={0} />
                </Form.Item>
                <Form.Item
                  name="electricAmount"
                  label="电费（分）"
                >
                  <InputNumber style={{ width: '100%' }} placeholder="请输入电费（分）" min={0} />
                </Form.Item>
                <Form.Item
                  name="otherAmount"
                  label="其他费用（分）"
                >
                  <InputNumber style={{ width: '100%' }} placeholder="请输入其他费用（分）" min={0} />
                </Form.Item>
              </div>
            </>
          )}

          <Form.Item
            name="paidAt"
            label="到账时间"
          >
            <DatePicker style={{ width: '100%' }} placeholder="请选择到账时间" />
          </Form.Item>
          <Form.Item
            name="note"
            label="备注"
          >
            <Input.TextArea placeholder="请输入备注" rows={3} />
          </Form.Item>
        </Form>
      </Modal>

      {/* 详情模态框 */}
      <Modal
        title="账单详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={[
          <Button key="close" onClick={() => setDetailModalVisible(false)}>
            关闭
          </Button>,
        ]}
        width={600}
      >
        {viewingBill && (
          <div>
            <h3 style={{ marginBottom: 16, paddingBottom: 12, borderBottom: '1px solid #f0f0f0' }}>
              基本信息
            </h3>
            <table style={{ width: '100%', marginBottom: 20 }}>
              <tbody>
                <tr>
                  <td style={{ width: 120, padding: '8px 0', fontWeight: 'bold' }}>ID:</td>
                  <td style={{ padding: '8px 0' }}>{viewingBill.id}</td>
                </tr>
                <tr>
                  <td style={{ padding: '8px 0', fontWeight: 'bold' }}>类型:</td>
                  <td style={{ padding: '8px 0' }}>
                    <Tag color={typeColorMap[viewingBill.type] || 'default'}>
                      {typeMap[viewingBill.type] || viewingBill.type}
                    </Tag>
                  </td>
                </tr>
                <tr>
                  <td style={{ padding: '8px 0', fontWeight: 'bold' }}>状态:</td>
                  <td style={{ padding: '8px 0' }}>
                    <Tag color={statusColorMap[viewingBill.status] || 'default'}>
                      {statusMap[viewingBill.status] || viewingBill.status}
                    </Tag>
                  </td>
                </tr>
                <tr>
                  <td style={{ padding: '8px 0', fontWeight: 'bold' }}>租约ID:</td>
                  <td style={{ padding: '8px 0' }}>{viewingBill.leaseId}</td>
                </tr>
                <tr>
                  <td style={{ padding: '8px 0', fontWeight: 'bold' }}>创建时间:</td>
                  <td style={{ padding: '8px 0' }}>{viewingBill.createdAt}</td>
                </tr>
                {viewingBill.paidAt && (
                  <tr>
                    <td style={{ padding: '8px 0', fontWeight: 'bold' }}>到账时间:</td>
                    <td style={{ padding: '8px 0' }}>{viewingBill.paidAt}</td>
                  </tr>
                )}
              </tbody>
            </table>

            <h3 style={{ marginBottom: 16, paddingBottom: 12, borderBottom: '1px solid #f0f0f0' }}>
              金额明细
            </h3>
            <div style={{ marginBottom: 20 }}>
              {viewingBill.type === 'checkout' ? (
                // 退租结算账单显示方式 - 所有项目都显示，即使是0
                <div>
                  <div style={{ padding: '8px 0', display: 'flex', justifyContent: 'space-between', borderBottom: '1px solid #f0f0f0' }}>
                    <span style={{ color: '#1890ff' }}>退还租金:</span>
                    <span style={{ color: '#1890ff', fontWeight: 'bold' }}>
                      {formatAmount(Math.abs(viewingBill.rentAmount || 0))}
                    </span>
                  </div>
                  <div style={{ padding: '8px 0', display: 'flex', justifyContent: 'space-between', borderBottom: '1px solid #f0f0f0' }}>
                    <span style={{ color: '#1890ff' }}>退还押金:</span>
                    <span style={{ color: '#1890ff', fontWeight: 'bold' }}>
                      {formatAmount(viewingBill.refundDepositAmount || 0)}
                    </span>
                  </div>
                  <div style={{ padding: '8px 0', display: 'flex', justifyContent: 'space-between', borderBottom: '1px solid #f0f0f0' }}>
                    <span style={{ color: '#faad14' }}>水费:</span>
                    <span style={{ color: '#faad14', fontWeight: 'bold' }}>
                      {formatAmount(viewingBill.waterAmount || 0)}
                    </span>
                  </div>
                  <div style={{ padding: '8px 0', display: 'flex', justifyContent: 'space-between', borderBottom: '1px solid #f0f0f0' }}>
                    <span style={{ color: '#faad14' }}>电费:</span>
                    <span style={{ color: '#faad14', fontWeight: 'bold' }}>
                      {formatAmount(viewingBill.electricAmount || 0)}
                    </span>
                  </div>
                  <div style={{ padding: '8px 0', display: 'flex', justifyContent: 'space-between', borderBottom: '1px solid #f0f0f0' }}>
                    <span style={{ color: '#faad14' }}>其他费用:</span>
                    <span style={{ color: '#faad14', fontWeight: 'bold' }}>
                      {formatAmount(viewingBill.otherAmount || 0)}
                    </span>
                  </div>
                  <div style={{ padding: '12px 0', display: 'flex', justifyContent: 'space-between', background: '#f5f5f5', marginTop: 8, paddingLeft: 8, paddingRight: 8, borderRadius: 4 }}>
                    <span style={{ fontWeight: 'bold', fontSize: '16px' }}>合计:</span>
                    <span style={{ fontWeight: 'bold', fontSize: '16px' }}>
                      {formatAmount(viewingBill.amount)}
                    </span>
                  </div>
                </div>
              ) : (
                // 普通账单显示方式
                <div>
                  {(viewingBill.rentAmount || 0) > 0 && (
                    <div style={{ padding: '8px 0', display: 'flex', justifyContent: 'space-between', borderBottom: '1px solid #f0f0f0' }}>
                      <span>租金:</span>
                      <span style={{ fontWeight: 'bold' }}>
                        {formatAmount(viewingBill.rentAmount || 0)}
                      </span>
                    </div>
                  )}
                  {(viewingBill.waterAmount || 0) > 0 && (
                    <div style={{ padding: '8px 0', display: 'flex', justifyContent: 'space-between', borderBottom: '1px solid #f0f0f0' }}>
                      <span>水费:</span>
                      <span style={{ fontWeight: 'bold' }}>
                        {formatAmount(viewingBill.waterAmount || 0)}
                      </span>
                    </div>
                  )}
                  {(viewingBill.electricAmount || 0) > 0 && (
                    <div style={{ padding: '8px 0', display: 'flex', justifyContent: 'space-between', borderBottom: '1px solid #f0f0f0' }}>
                      <span>电费:</span>
                      <span style={{ fontWeight: 'bold' }}>
                        {formatAmount(viewingBill.electricAmount || 0)}
                      </span>
                    </div>
                  )}
                  {(viewingBill.otherAmount || 0) > 0 && (
                    <div style={{ padding: '8px 0', display: 'flex', justifyContent: 'space-between', borderBottom: '1px solid #f0f0f0' }}>
                      <span>其他:</span>
                      <span style={{ fontWeight: 'bold' }}>
                        {formatAmount(viewingBill.otherAmount || 0)}
                      </span>
                    </div>
                  )}
                  <div style={{ padding: '12px 0', display: 'flex', justifyContent: 'space-between', background: '#f5f5f5', marginTop: 8, paddingLeft: 8, paddingRight: 8, borderRadius: 4 }}>
                    <span style={{ fontWeight: 'bold', fontSize: '16px' }}>合计:</span>
                    <span style={{ fontWeight: 'bold', fontSize: '16px' }}>
                      {formatAmount(viewingBill.amount)}
                    </span>
                  </div>
                </div>
              )}
            </div>

            {viewingBill.note && (
              <>
                <h3 style={{ marginBottom: 16, paddingBottom: 12, borderBottom: '1px solid #f0f0f0' }}>
                  备注
                </h3>
                <p style={{ margin: 0, padding: '8px 0', background: '#f5f5f5', paddingLeft: 12, paddingRight: 12, borderRadius: 4 }}>
                  {viewingBill.note}
                </p>
              </>
            )}
          </div>
        )}
      </Modal>

      <OperationLogModal
        visible={operationLogVisible}
        aggregateId={currentBill?.id || ''}
        domainType="bill"
        title={`账单操作日志 - ${currentBill?.id || ''}`}
        onCancel={() => setOperationLogVisible(false)}
      />
    </div>
  );
};

export default Bills;
