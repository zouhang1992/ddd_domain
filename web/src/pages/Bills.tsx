import React, { useState, useEffect, useMemo } from 'react';
import { Table, Button, Modal, Form, Input, message, Popconfirm, Select, DatePicker, InputNumber, Tag, Pagination } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, PrinterOutlined, HistoryOutlined, SearchOutlined, ReloadOutlined, EyeOutlined } from '@ant-design/icons';
import type { Bill, Lease, Room, Location, Landlord } from '../types/api';
import { billApi, type BillsQueryResult } from '../api/bill';
import { leaseApi } from '../api/lease';
import { roomApi } from '../api/room';
import { locationApi } from '../api/location';
import { landlordApi } from '../api/landlord';
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
  const [allBills, setAllBills] = useState<Bill[]>([]);
  const [displayBills, setDisplayBills] = useState<Bill[]>([]);
  const [leases, setLeases] = useState<Lease[]>([]);
  const [rooms, setRooms] = useState<Room[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);
  const [landlords, setLandlords] = useState<Landlord[]>([]);
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
  const [queryLocationId, setQueryLocationId] = useState<string>();
  const [queryRoomId, setQueryRoomId] = useState<string>();
  const [queryType, setQueryType] = useState<string>();
  const [queryStatus, setQueryStatus] = useState<string>();

  const fetchBills = async () => {
    setLoading(true);
    try {
      const data: BillsQueryResult = await billApi.list({ limit: 1000 });
      setAllBills(data.items || []);
      setTotal(data.total);
    } catch {
      message.error('获取账单列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchLeases = async () => {
    try {
      const data = await leaseApi.list({ limit: 1000 });
      setLeases(data.items || []);
    } catch {
      message.error('获取租约列表失败');
    }
  };

  const fetchRooms = async () => {
    try {
      const data = await roomApi.list();
      setRooms(data.items || []);
    } catch {
      message.error('获取房间列表失败');
    }
  };

  const fetchLocations = async () => {
    try {
      const data = await locationApi.list();
      setLocations(data.items || []);
    } catch {
      message.error('获取位置列表失败');
    }
  };

  const fetchLandlords = async () => {
    try {
      const data = await landlordApi.list();
      setLandlords(data.items || []);
    } catch {
      message.error('获取房东列表失败');
    }
  };

  useEffect(() => {
    fetchBills();
    fetchLeases();
    fetchRooms();
    fetchLocations();
    fetchLandlords();
  }, []);

  // 筛选用于查询的房间
  const queryRooms = useMemo(() => {
    if (queryLocationId) {
      return rooms.filter(room => room.locationId === queryLocationId);
    }
    return rooms;
  }, [rooms, queryLocationId]);

  // 应用筛选和分页
  useEffect(() => {
    let filtered = [...allBills];

    // 位置筛选 - 通过房间关联
    if (queryLocationId) {
      const roomIds = rooms.filter(r => r.locationId === queryLocationId).map(r => r.id);
      const leaseIds = leases.filter(l => roomIds.includes(l.roomId)).map(l => l.id);
      filtered = filtered.filter(b => leaseIds.includes(b.leaseId));
    }

    // 房间筛选 - 通过租约关联
    if (queryRoomId) {
      const leaseIds = leases.filter(l => l.roomId === queryRoomId).map(l => l.id);
      filtered = filtered.filter(b => leaseIds.includes(b.leaseId));
    }

    // 类型筛选
    if (queryType) {
      filtered = filtered.filter(b => b.type === queryType);
    }

    // 状态筛选
    if (queryStatus) {
      filtered = filtered.filter(b => b.status === queryStatus);
    }

    // 分页
    const start = (page - 1) * pageSize;
    const end = start + pageSize;
    setDisplayBills(filtered.slice(start, end));
    setTotal(filtered.length);
  }, [allBills, leases, rooms, page, pageSize, queryLocationId, queryRoomId, queryType, queryStatus]);

  const handleQuery = async () => {
    const values = await queryForm.validateFields();
    setQueryType(values.type);
    setQueryStatus(values.status);
    setPage(1);
  };

  const handleReset = () => {
    queryForm.resetFields();
    setQueryLocationId(undefined);
    setQueryRoomId(undefined);
    setQueryType(undefined);
    setQueryStatus(undefined);
    setPage(1);
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
      billStart: bill.billStart ? dayjs(bill.billStart) : null,
      billEnd: bill.billEnd ? dayjs(bill.billEnd) : null,
      dueDate: bill.dueDate ? dayjs(bill.dueDate) : null,
      refundDepositAmount: bill.refundDepositAmount || 0,
    });
    setModalVisible(true);
  };

  const handlePrint = async (bill: Bill) => {
    try {
      await billApi.printReceipt(bill.id);
      message.success('打印成功');
    } catch {
      message.error('打印失败');
    }
  };

  const handleConfirmArrival = async (bill: Bill) => {
    try {
      await billApi.confirmArrival(bill.id);
      message.success('到账确认成功');
      fetchBills();
    } catch {
      message.error('到账确认失败');
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await billApi.delete(id);
      message.success('删除成功');
      fetchBills();
    } catch {
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
        billStart: values.billStart ? values.billStart.format('YYYY-MM-DD') : '',
        billEnd: values.billEnd ? values.billEnd.format('YYYY-MM-DD') : '',
        dueDate: values.dueDate ? values.dueDate.format('YYYY-MM-DD') : '',
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
    } catch {
      message.error('操作失败');
    }
  };

  const formatAmount = (amount: number) => {
    return `¥${(amount / 100).toFixed(2)}`;
  };

  const columns = [
    {
      title: '位置',
      key: 'location',
      render: (_: any, record: Bill) => {
        const lease = leases.find(l => l.id === record.leaseId);
        if (!lease) return '-';
        const room = rooms.find(r => r.id === lease.roomId);
        const location = room ? locations.find(l => l.id === room.locationId) : null;
        return location?.shortName || '-';
      },
    },
    {
      title: '房间',
      key: 'room',
      render: (_: unknown, record: Bill) => {
        const lease = leases.find(l => l.id === record.leaseId);
        if (!lease) return '-';
        const room = rooms.find(r => r.id === lease.roomId);
        if (!room) return '-';
        const location = locations.find(l => l.id === room.locationId);
        return (
          <span>
            <Tag color="blue" style={{ marginRight: 8 }}>
              {location?.shortName || '未知位置'}
            </Tag>
            {room.roomNumber}
          </span>
        );
      },
    },
    {
      title: '地址',
      key: 'address',
      render: (_: unknown, record: Bill) => {
        const lease = leases.find(l => l.id === record.leaseId);
        if (!lease) return '-';
        const room = rooms.find(r => r.id === lease.roomId);
        if (!room) return '-';
        const location = locations.find(l => l.id === room.locationId);
        return location?.detail || location?.shortName || '-';
      },
    },
    {
      title: '房东',
      key: 'landlord',
      render: (_: unknown, record: Bill) => {
        const lease = leases.find(l => l.id === record.leaseId);
        if (!lease) return '-';
        return landlords.find(l => l.id === lease.landlordId)?.name || '-';
      },
    },
    {
      title: '租期',
      key: 'leasePeriod',
      width: 180,
      render: (_: any, record: Bill) => {
        const lease = leases.find(l => l.id === record.leaseId);
        if (!lease) return '-';
        return `${dayjs(lease.startDate).format('YYYY-MM-DD')} ~ ${dayjs(lease.endDate).format('YYYY-MM-DD')}`;
      },
    },
    {
      title: '计费周期',
      key: 'billPeriod',
      width: 180,
      render: (_: any, record: Bill) => {
        if (record.billStart && record.billEnd) {
          return `${dayjs(record.billStart).format('YYYY-MM-DD')} ~ ${dayjs(record.billEnd).format('YYYY-MM-DD')}`;
        }
        return '-';
      },
    },
    {
      title: '租户',
      key: 'tenant',
      width: 100,
      render: (_: any, record: Bill) => {
        const lease = leases.find(l => l.id === record.leaseId);
        return lease?.tenantName || '-';
      },
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
        const rent = record.rentAmount || 0;
        const water = record.waterAmount || 0;
        const electric = record.electricAmount || 0;
        const other = record.otherAmount || 0;

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
        return formatAmount(record.amount);
      },
    },
    {
      title: '到账时间',
      dataIndex: 'paidAt',
      key: 'paidAt',
      render: (paidAt: string) => paidAt ? dayjs(paidAt).format('YYYY-MM-DD HH:mm:ss') : '-',
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (createdAt: string) => createdAt ? dayjs(createdAt).format('YYYY-MM-DD HH:mm:ss') : '-',
    },
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
        <Form.Item label="位置">
          <Select
            placeholder="请选择位置"
            style={{ width: 150 }}
            allowClear
            value={queryLocationId}
            onChange={(value) => { setQueryLocationId(value); setQueryRoomId(undefined); setPage(1); }}
          >
            {locations.map(location => (
              <Option key={location.id} value={location.id}>
                {location.shortName}
              </Option>
            ))}
          </Select>
        </Form.Item>
        <Form.Item label="房间">
          <Select
            placeholder="请选择房间"
            style={{ width: 150 }}
            allowClear
            value={queryRoomId}
            onChange={(value) => { setQueryRoomId(value); setPage(1); }}
          >
            {queryRooms.map(room => {
              const location = locations.find(l => l.id === room.locationId);
              return (
                <Option key={room.id} value={room.id}>
                  [{location?.shortName || '未知位置'}] {room.roomNumber}
                </Option>
              );
            })}
          </Select>
        </Form.Item>
        <Form.Item label="类型">
          <Select
            placeholder="请选择类型"
            style={{ width: 120 }}
            allowClear
            value={queryType}
            onChange={(value) => { setQueryType(value); setPage(1); }}
          >
            <Option value="rent">租金</Option>
            <Option value="charge">收账</Option>
            <Option value="checkout">退租结算</Option>
          </Select>
        </Form.Item>
        <Form.Item label="状态">
          <Select
            placeholder="请选择状态"
            style={{ width: 120 }}
            allowClear
            value={queryStatus}
            onChange={(value) => { setQueryStatus(value); setPage(1); }}
          >
            <Option value="pending">待到账</Option>
            <Option value="paid">已到账</Option>
          </Select>
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
        dataSource={displayBills}
        rowKey="id"
        loading={loading}
        pagination={false}
        scroll={{ x: 2800 }}
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
          {!editingBill ? (
            <>
              <Form.Item
                name="leaseId"
                label="租约"
                rules={[{ required: true, message: '请选择租约' }]}
              >
                <Select
                  placeholder="请选择租约"
                  showSearch
                  optionFilterProp="children"
                  onChange={async (leaseId) => {
                    if (leaseId) {
                      try {
                        const data = await billApi.getNextBillPeriod(leaseId);
                        form.setFieldValue('billStart', dayjs(data.billStart));
                        form.setFieldValue('billEnd', null);
                      } catch {
                        message.error('获取默认计费开始时间失败');
                      }
                    }
                  }}
                >
                  {leases
                    .filter(lease => lease.status !== 'expired' && lease.status !== 'checkout')
                    .map(lease => {
                      const room = rooms.find(r => r.id === lease.roomId);
                      const location = room ? locations.find(l => l.id === room.locationId) : null;
                      return (
                        <Option key={lease.id} value={lease.id}>
                          <div>
                            <div style={{ fontWeight: 500 }}>
                              {location?.shortName || '未知位置'} - {room?.roomNumber || lease.roomId}
                            </div>
                            <div style={{ fontSize: '12px', color: '#666' }}>
                              租户: {lease?.tenantName} | {lease?.tenantPhone} | 租期: {lease.startDate} ~ {lease.endDate}
                            </div>
                          </div>
                        </Option>
                      );
                    })}
                </Select>
              </Form.Item>
              {/* 显示选中租约的位置和房间信息 */}
              {form.getFieldValue('leaseId') && (
                <div style={{ marginBottom: 16, padding: 12, background: '#f5f5f5', borderRadius: 6 }}>
                  {(() => {
                    const leaseId = form.getFieldValue('leaseId');
                    const lease = leases.find(l => l.id === leaseId);
                    if (!lease) return null;
                    const room = rooms.find(r => r.id === lease.roomId);
                    const location = room ? locations.find(l => l.id === room.locationId) : null;
                    return (
                      <>
                        <div style={{ marginBottom: 8, fontWeight: 500 }}>租约信息</div>
                        <div style={{ fontSize: '13px' }}>
                          <div><strong>位置:</strong> {location?.shortName || '-'}</div>
                          <div><strong>房间:</strong> {room?.roomNumber || '-'}</div>
                          <div><strong>地址:</strong> {location?.detail || '-'}</div>
                          <div><strong>租户:</strong> {lease?.tenantName} ({lease?.tenantPhone})</div>
                          <div><strong>租期:</strong> {lease.startDate} ~ {lease.endDate}</div>
                          <div><strong>状态:</strong> {lease.status}</div>
                        </div>
                      </>
                    );
                  })()}
                </div>
              )}
            </>
          ) : (
            /* 编辑账单时显示租约信息 */
            <div style={{ marginBottom: 16, padding: 12, background: '#f5f5f5', borderRadius: 6 }}>
              {(() => {
                const lease = leases.find(l => l.id === editingBill.leaseId);
                if (!lease) return null;
                const room = rooms.find(r => r.id === lease.roomId);
                const location = room ? locations.find(l => l.id === room.locationId) : null;
                return (
                  <>
                    <div style={{ marginBottom: 8, fontWeight: 500 }}>租约信息</div>
                    <div style={{ fontSize: '13px' }}>
                      <div><strong>位置:</strong> {location?.shortName || '-'}</div>
                      <div><strong>房间:</strong> {room?.roomNumber || '-'}</div>
                      <div><strong>地址:</strong> {location?.detail || '-'}</div>
                      <div><strong>租户:</strong> {lease?.tenantName} ({lease?.tenantPhone})</div>
                      <div><strong>租期:</strong> {lease.startDate} ~ {lease.endDate}</div>
                      <div><strong>状态:</strong> {lease.status}</div>
                    </div>
                  </>
                );
              })()}
            </div>
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
              <Select placeholder="请选择类型">
                <Option value="rent">租金</Option>
              </Select>
            </Form.Item>
          )}

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
            name="billStart"
            label="计费开始日期"
            rules={[{ required: true, message: '请选择计费开始日期' }]}
          >
            <DatePicker style={{ width: '100%' }} placeholder="请选择计费开始日期" />
          </Form.Item>
          <Form.Item
            name="billEnd"
            label="计费结束日期"
            rules={[{ required: true, message: '请选择计费结束日期' }]}
          >
            <DatePicker style={{ width: '100%' }} placeholder="请选择计费结束日期" />
          </Form.Item>
          <Form.Item
            name="dueDate"
            label="付款截止日期"
            rules={[{ required: true, message: '请选择付款截止日期' }]}
          >
            <DatePicker style={{ width: '100%' }} placeholder="请选择付款截止日期" />
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
                  <td style={{ padding: '8px 0' }}>{dayjs(viewingBill.createdAt).format('YYYY-MM-DD HH:mm:ss')}</td>
                </tr>
                {viewingBill.paidAt && (
                  <tr>
                    <td style={{ padding: '8px 0', fontWeight: 'bold' }}>到账时间:</td>
                    <td style={{ padding: '8px 0' }}>{dayjs(viewingBill.paidAt).format('YYYY-MM-DD HH:mm:ss')}</td>
                  </tr>
                )}
              </tbody>
            </table>

            <h3 style={{ marginBottom: 16, paddingBottom: 12, borderBottom: '1px solid #f0f0f0' }}>
              金额明细
            </h3>
            <div style={{ marginBottom: 20 }}>
              {viewingBill.type === 'checkout' ? (
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
