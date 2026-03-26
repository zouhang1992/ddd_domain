import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, message, Popconfirm, Select, DatePicker, InputNumber, Tag, Pagination } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, PrinterOutlined, HistoryOutlined, SearchOutlined, ReloadOutlined } from '@ant-design/icons';
import type { Bill, Lease } from '../types/api';
import { billApi, type BillQueryParams, type BillsQueryResult } from '../api/bill';
import { leaseApi } from '../api/lease';
import dayjs from 'dayjs';
import OperationLogModal from '../components/OperationLogModal';

const { Option } = Select;

const typeMap: Record<string, string> = {
  charge: '收账',
  checkout: '退租结算',
};

const typeColorMap: Record<string, string> = {
  charge: 'success',
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
  const [editingBill, setEditingBill] = useState<Bill | null>(null);
  const [form] = Form.useForm();
  const [operationLogVisible, setOperationLogVisible] = useState(false);
  const [currentBill, setCurrentBill] = useState<Bill | null>(null);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [queryForm] = Form.useForm();

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
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (bill: Bill) => {
    setEditingBill(bill);
    form.setFieldsValue({
      ...bill,
      paidAt: bill.paidAt ? dayjs(bill.paidAt) : null,
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
      title: '总金额',
      dataIndex: 'amount',
      key: 'amount',
      render: formatAmount,
    },
    {
      title: '租金',
      dataIndex: 'rentAmount',
      key: 'rentAmount',
      render: formatAmount,
    },
    {
      title: '水费',
      dataIndex: 'waterAmount',
      key: 'waterAmount',
      render: formatAmount,
    },
    {
      title: '电费',
      dataIndex: 'electricAmount',
      key: 'electricAmount',
      render: formatAmount,
    },
    {
      title: '其他',
      dataIndex: 'otherAmount',
      key: 'otherAmount',
      render: formatAmount,
    },
    { title: '到账时间', dataIndex: 'paidAt', key: 'paidAt' },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: Bill) => (
        <div>
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
            <Option value="charge">收账</Option>
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
          <Form.Item
            name="type"
            label="类型"
            rules={[{ required: true, message: '请选择类型' }]}
          >
            <Select placeholder="请选择类型">
              <Option value="charge">收账</Option>
              <Option value="checkout">退租结算</Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="amount"
            label="总金额（分）"
            rules={[{ required: true, message: '请输入总金额' }]}
          >
            <InputNumber style={{ width: '100%' }} placeholder="请输入总金额（分）" />
          </Form.Item>
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
