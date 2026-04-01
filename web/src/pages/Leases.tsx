import React, { useState, useEffect, useMemo } from 'react';
import { Table, Button, Modal, Form, Input, message, Popconfirm, Select, DatePicker, InputNumber, Tag, Space, Pagination } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, CheckOutlined, ReloadOutlined, PrinterOutlined, ThunderboltOutlined, HistoryOutlined, SearchOutlined } from '@ant-design/icons';
import type { Lease, Room, Landlord, Location } from '../types/api';
import { leaseApi, type LeasesQueryResult } from '../api/lease';
import { roomApi } from '../api/room';
import { landlordApi } from '../api/landlord';
import { locationApi } from '../api/location';
import dayjs from 'dayjs';
import OperationLogModal from '../components/OperationLogModal';

const { Option } = Select;

const statusMap: Record<string, string> = {
  pending: '待生效',
  active: '生效中',
  expired: '已过期',
  checkout: '已退租',
};

const statusColorMap: Record<string, string> = {
  pending: 'default',
  active: 'success',
  expired: 'warning',
  checkout: 'error',
};

const Leases: React.FC = () => {
  const [allLeases, setAllLeases] = useState<Lease[]>([]);
  const [displayLeases, setDisplayLeases] = useState<Lease[]>([]);
  const [rooms, setRooms] = useState<Room[]>([]);
  const [landlords, setLandlords] = useState<Landlord[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [renewModalVisible, setRenewModalVisible] = useState(false);
  const [checkoutModalVisible, setCheckoutModalVisible] = useState(false);
  const [editingLease, setEditingLease] = useState<Lease | null>(null);
  const [checkoutLease, setCheckoutLease] = useState<Lease | null>(null);
  const [form] = Form.useForm();
  const [renewForm] = Form.useForm();
  const [checkoutForm] = Form.useForm();
  const [selectedLocationId, setSelectedLocationId] = useState<string>('');
  const [operationLogVisible, setOperationLogVisible] = useState(false);
  const [currentLease, setCurrentLease] = useState<Lease | null>(null);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [queryForm] = Form.useForm();
  const [queryLocationId, setQueryLocationId] = useState<string>();
  const [queryRoomId, setQueryRoomId] = useState<string>();
  const [queryTenantName, setQueryTenantName] = useState<string>();
  const [queryTenantPhone, setQueryTenantPhone] = useState<string>();
  const [queryStatus, setQueryStatus] = useState<string>();

  const fetchLeases = async () => {
    setLoading(true);
    try {
      const data: LeasesQueryResult = await leaseApi.list({ limit: 1000 });
      setAllLeases(data.items || []);
      setTotal(data.total);
    } catch {
      message.error('获取租约列表失败');
    } finally {
      setLoading(false);
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

  const fetchLandlords = async () => {
    try {
      const data = await landlordApi.list();
      setLandlords(data.items || []);
    } catch {
      message.error('获取房东列表失败');
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

  useEffect(() => {
    fetchLeases();
    fetchRooms();
    fetchLandlords();
    fetchLocations();
  }, []);

  // 应用筛选和分页
  useEffect(() => {
    let filtered = [...allLeases];

    // 位置筛选
    if (queryLocationId) {
      const roomIds = rooms.filter(r => r.locationId === queryLocationId).map(r => r.id);
      filtered = filtered.filter(l => roomIds.includes(l.roomId));
    }

    // 房间筛选
    if (queryRoomId) {
      filtered = filtered.filter(l => l.roomId === queryRoomId);
    }

    // 租户姓名筛选
    if (queryTenantName) {
      filtered = filtered.filter(l =>
        l.tenantName?.toLowerCase().includes(queryTenantName.toLowerCase())
      );
    }

    // 租户电话筛选
    if (queryTenantPhone) {
      filtered = filtered.filter(l =>
        l.tenantPhone?.includes(queryTenantPhone)
      );
    }

    // 状态筛选
    if (queryStatus) {
      filtered = filtered.filter(l => l.status === queryStatus);
    }

    // 分页
    const start = (page - 1) * pageSize;
    const end = start + pageSize;
    setDisplayLeases(filtered.slice(start, end));
    setTotal(filtered.length);
  }, [allLeases, rooms, page, pageSize, queryLocationId, queryRoomId, queryTenantName, queryTenantPhone, queryStatus]);

  const handleQuery = async () => {
    const values = await queryForm.validateFields();
    setQueryTenantName(values.tenantName);
    setQueryTenantPhone(values.tenantPhone);
    setQueryStatus(values.status);
    setPage(1);
  };

  const handleReset = () => {
    queryForm.resetFields();
    setQueryLocationId(undefined);
    setQueryRoomId(undefined);
    setQueryTenantName(undefined);
    setQueryTenantPhone(undefined);
    setQueryStatus(undefined);
    setPage(1);
  };

  const handlePageChange = (pageNum: number, pageSizeNum: number) => {
    setPage(pageNum);
    setPageSize(pageSizeNum);
  };

  const handleCreate = () => {
    setEditingLease(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (lease: Lease) => {
    setEditingLease(lease);
    form.setFieldsValue({
      ...lease,
      startDate: dayjs(lease.startDate),
      endDate: dayjs(lease.endDate),
    });
    setModalVisible(true);
  };

  const handleRenew = (lease: Lease) => {
    setEditingLease(lease);
    renewForm.resetFields();
    renewForm.setFieldsValue({
      newStartDate: dayjs(lease.endDate).add(1, 'day'),
      newRentAmount: lease.rentAmount,
    });
    setRenewModalVisible(true);
  };

  const handleCheckout = (lease: Lease) => {
    setCheckoutLease(lease);
    checkoutForm.resetFields();
    checkoutForm.setFieldsValue({
      refundRentAmount: 0,
      refundDepositAmount: lease.depositAmount || 0,
      waterAmount: 0,
      electricAmount: 0,
      otherAmount: 0,
      note: '',
    });
    setCheckoutModalVisible(true);
  };

  const handleCheckoutSubmit = async () => {
    if (!checkoutLease) return;
    try {
      const values = await checkoutForm.validateFields();
      await leaseApi.checkoutWithBills(checkoutLease.id, values);
      message.success('退租成功');
      setCheckoutModalVisible(false);
      fetchLeases();
    } catch {
      message.error('退租失败');
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await leaseApi.delete(id);
      message.success('删除成功');
      fetchLeases();
    } catch {
      message.error('删除失败');
    }
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      const formattedValues = {
        ...values,
        startDate: values.startDate.format('YYYY-MM-DD'),
        endDate: values.endDate.format('YYYY-MM-DD'),
      };
      if (editingLease) {
        await leaseApi.update(editingLease.id, formattedValues);
        message.success('更新成功');
      } else {
        await leaseApi.create(formattedValues);
        message.success('创建成功');
      }
      setModalVisible(false);
      fetchLeases();
    } catch {
      message.error('操作失败');
    }
  };

  const handleRenewModalOk = async () => {
    if (!editingLease) return;
    try {
      const values = await renewForm.validateFields();
      const formattedValues = {
        ...values,
        newStartDate: values.newStartDate.format('YYYY-MM-DD'),
        newEndDate: values.newEndDate.format('YYYY-MM-DD'),
      };
      await leaseApi.renew(editingLease.id, formattedValues);
      message.success('续租成功');
      setRenewModalVisible(false);
      fetchLeases();
    } catch {
      message.error('续租失败');
    }
  };

  const formatAmount = (amount: number) => {
    return `¥${(amount / 100).toFixed(2)}`;
  };

  // 根据选择的位置筛选房间，同时过滤掉已出租的房间
  const filteredRooms = useMemo(() => {
    let result = rooms;
    if (selectedLocationId) {
      result = result.filter(room => room.locationId === selectedLocationId);
    }
    // 只显示可出租的房间
    return result.filter(room => room.status === 'available');
  }, [rooms, selectedLocationId]);

  // 筛选用于查询的房间
  const queryRooms = useMemo(() => {
    if (queryLocationId) {
      return rooms.filter(room => room.locationId === queryLocationId);
    }
    return rooms;
  }, [rooms, queryLocationId]);

  // 处理位置选择变化
  const handleLocationChange = (value: string) => {
    setSelectedLocationId(value);
    form.setFieldValue('roomId', undefined);
  };

  // 处理查询位置变化
  const handleQueryLocationChange = (value: string | undefined) => {
    setQueryLocationId(value);
    setQueryRoomId(undefined);
    setPage(1);
  };

  const handlePrintContract = async (lease: Lease) => {
    try {
      await leaseApi.printContract(lease.id);
      message.success('合同下载成功');
    } catch {
      message.error('合同下载失败');
    }
  };

  const handleActivateLease = async (lease: Lease) => {
    try {
      await leaseApi.activate(lease.id);
      message.success('租约生效成功');
      fetchLeases();
    } catch {
      message.error('租约生效失败');
    }
  };

  const handleViewOperationLogs = (lease: Lease) => {
    setCurrentLease(lease);
    setOperationLogVisible(true);
  };

  const columns = [
    {
      title: '位置',
      key: 'location',
      render: (_: any, record: Lease) => {
        const room = rooms.find(r => r.id === record.roomId);
        const location = room ? locations.find(l => l.id === room.locationId) : null;
        return location?.shortName || '-';
      },
    },
    {
      title: '房间',
      dataIndex: 'roomId',
      key: 'roomId',
      render: (roomId: string) => {
        const room = rooms.find(r => r.id === roomId);
        if (!room) return roomId;
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
      title: '房东',
      key: 'landlord',
      render: (_: unknown, record: Lease) => {
        const landlord = landlords.find(l => l.id === record.landlordId);
        return landlord?.name || '-';
      },
    },
    { title: '租户姓名', dataIndex: 'tenantName', key: 'tenantName' },
    { title: '租户电话', dataIndex: 'tenantPhone', key: 'tenantPhone' },
    {
      title: '开始日期',
      dataIndex: 'startDate',
      key: 'startDate',
      render: (date: string) => dayjs(date).format('YYYY-MM-DD'),
    },
    {
      title: '结束日期',
      dataIndex: 'endDate',
      key: 'endDate',
      render: (date: string) => dayjs(date).format('YYYY-MM-DD'),
    },
    {
      title: '租金',
      dataIndex: 'rentAmount',
      key: 'rentAmount',
      render: formatAmount,
    },
    {
      title: '押金金额',
      dataIndex: 'depositAmount',
      key: 'depositAmount',
      render: (deposit: number) => formatAmount(deposit || 0),
    },
    {
      title: '押金状态',
      key: 'depositStatus',
      render: (_: any, record: Lease) => {
        if (record.status === 'checkout') {
          return <Tag color="default">已退还</Tag>;
        }
        return <Tag color="success">已收取</Tag>;
      },
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
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (createdAt: string) => createdAt ? dayjs(createdAt).format('YYYY-MM-DD HH:mm:ss') : '-',
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: Lease) => (
        <Space size="small">
          <Button
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
            size="small"
          >
            编辑
          </Button>
          {record.status === 'active' && (
            <Button
              icon={<ReloadOutlined />}
              onClick={() => handleRenew(record)}
              size="small"
            >
              续租
            </Button>
          )}
          {record.status === 'pending' && (
            <Button
              icon={<ThunderboltOutlined />}
              onClick={() => handleActivateLease(record)}
              type="primary"
              size="small"
            >
              生效
            </Button>
          )}
          {record.status === 'active' && (
            <Popconfirm
              title="确定退租?"
              onConfirm={() => handleCheckout(record)}
              okText="确定"
              cancelText="取消"
            >
              <Button
                icon={<CheckOutlined />}
                type="default"
                size="small"
              >
                退租
              </Button>
            </Popconfirm>
          )}
          <Button
            icon={<PrinterOutlined />}
            onClick={() => handlePrintContract(record)}
            size="small"
          >
            打印合同
          </Button>
          <Button
            icon={<HistoryOutlined />}
            onClick={() => handleViewOperationLogs(record)}
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
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h1>租约管理</h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新增租约
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
            onChange={handleQueryLocationChange}
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
        <Form.Item name="tenantName" label="租户姓名">
          <Input
            placeholder="请输入租户姓名"
            style={{ width: 120 }}
            value={queryTenantName}
            onChange={(e) => setQueryTenantName(e.target.value)}
          />
        </Form.Item>
        <Form.Item name="tenantPhone" label="租户电话">
          <Input
            placeholder="请输入租户电话"
            style={{ width: 120 }}
            value={queryTenantPhone}
            onChange={(e) => setQueryTenantPhone(e.target.value)}
          />
        </Form.Item>
        <Form.Item label="状态">
          <Select
            placeholder="请选择状态"
            style={{ width: 120 }}
            allowClear
            value={queryStatus}
            onChange={(value) => { setQueryStatus(value); setPage(1); }}
          >
            <Option value="pending">待生效</Option>
            <Option value="active">生效中</Option>
            <Option value="expired">已过期</Option>
            <Option value="checkout">已退租</Option>
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
        dataSource={displayLeases}
        rowKey="id"
        loading={loading}
        pagination={false}
        scroll={{ x: 2000 }}
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
        title={editingLease ? '编辑租约' : '新增租约'}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={() => setModalVisible(false)}
        width={600}
      >
        <Form form={form} layout="vertical">
          {!editingLease && (
            <>
              <Form.Item
                name="locationId"
                label="位置"
                rules={[{ required: true, message: '请选择位置' }]}
              >
                <Select
                  placeholder="请选择位置"
                  onChange={handleLocationChange}
                >
                  {locations.map(location => (
                    <Option key={location.id} value={location.id}>
                      {location.shortName}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
              <Form.Item
                name="roomId"
                label="房间"
                rules={[{ required: true, message: '请选择房间' }]}
              >
                <Select
                  placeholder={
                    filteredRooms.length === 0
                      ? (selectedLocationId ? '该位置暂无可出租房间' : '请先选择位置，再选择房间')
                      : '请选择房间'
                  }
                >
                  {filteredRooms.map(room => {
                    const location = locations.find(l => l.id === room.locationId);
                    return (
                      <Option key={room.id} value={room.id}>
                        [{location?.shortName || '未知位置'}] {room.roomNumber}
                      </Option>
                    );
                  })}
                </Select>
              </Form.Item>
              <Form.Item
                name="landlordId"
                label="房东"
                rules={[{ required: true, message: '请选择房东' }]}
              >
                <Select placeholder="请选择房东">
                  {landlords.map(landlord => (
                    <Option key={landlord.id} value={landlord.id}>
                      {landlord.name}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
              <Form.Item
                name="depositAmount"
                label="押金（分）"
                rules={[{ required: true, message: '请输入押金' }]}
              >
                <InputNumber style={{ width: '100%' }} placeholder="请输入押金（分）" />
              </Form.Item>
              <Form.Item
                name="depositNote"
                label="押金备注"
              >
                <Input.TextArea placeholder="请输入押金备注" rows={2} />
              </Form.Item>
            </>
          )}
          <Form.Item
            name="tenantName"
            label="租户姓名"
            rules={[{ required: true, message: '请输入租户姓名' }]}
          >
            <Input placeholder="请输入租户姓名" />
          </Form.Item>
          <Form.Item
            name="tenantPhone"
            label="租户电话"
            rules={[{ required: true, message: '请输入租户电话' }]}
          >
            <Input placeholder="请输入租户电话" />
          </Form.Item>
          <Form.Item
            name="startDate"
            label="开始日期"
            rules={[{ required: true, message: '请选择开始日期' }]}
          >
            <DatePicker style={{ width: '100%' }} placeholder="请选择开始日期" />
          </Form.Item>
          <Form.Item
            name="endDate"
            label="结束日期"
            rules={[{ required: true, message: '请选择结束日期' }]}
          >
            <DatePicker style={{ width: '100%' }} placeholder="请选择结束日期" />
          </Form.Item>
          <Form.Item
            name="rentAmount"
            label="租金（分）"
            rules={[{ required: true, message: '请输入租金' }]}
          >
            <InputNumber style={{ width: '100%' }} placeholder="请输入租金（分）" />
          </Form.Item>
          <Form.Item
            name="note"
            label="备注"
          >
            <Input.TextArea placeholder="请输入备注" rows={3} />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="续租"
        open={renewModalVisible}
        onOk={handleRenewModalOk}
        onCancel={() => setRenewModalVisible(false)}
        width={500}
      >
        <Form form={renewForm} layout="vertical">
          <Form.Item
            name="newStartDate"
            label="新开始日期"
            rules={[{ required: true, message: '请选择新开始日期' }]}
          >
            <DatePicker style={{ width: '100%' }} placeholder="请选择新开始日期" />
          </Form.Item>
          <Form.Item
            name="newEndDate"
            label="新结束日期"
            rules={[{ required: true, message: '请选择新结束日期' }]}
          >
            <DatePicker style={{ width: '100%' }} placeholder="请选择新结束日期" />
          </Form.Item>
          <Form.Item
            name="newRentAmount"
            label="新租金（分）"
            rules={[{ required: true, message: '请输入新租金' }]}
          >
            <InputNumber style={{ width: '100%' }} placeholder="请输入新租金（分）" />
          </Form.Item>
          <Form.Item
            name="note"
            label="备注"
          >
            <Input.TextArea placeholder="请输入备注" rows={3} />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="退租结算"
        open={checkoutModalVisible}
        onOk={handleCheckoutSubmit}
        onCancel={() => setCheckoutModalVisible(false)}
        width={600}
      >
        <Form form={checkoutForm} layout="vertical">
          {/* 租约信息 */}
          <div style={{ marginBottom: 16, padding: 16, background: '#f5f5f5', borderRadius: 8 }}>
            <h4 style={{ margin: '0 0 12px 0', color: '#666' }}>租约信息</h4>
            <p style={{ margin: '0 0 8px 0' }}><strong>租户:</strong> {checkoutLease?.tenantName}</p>
            <p style={{ margin: 0 }}><strong>押金金额:</strong> {formatAmount(checkoutLease?.depositAmount || 0)}</p>
          </div>

          {/* 退还金额 */}
          <div style={{ marginBottom: 16, padding: 16, border: '1px solid #d9d9d9', borderRadius: 8 }}>
            <h4 style={{ margin: '0 0 16px 0', color: '#1890ff' }}>退还金额</h4>
            <Form.Item
              name="refundRentAmount"
              label="退还租金（分）"
              rules={[{ required: true, message: '请输入退还租金' }]}
            >
              <InputNumber
                style={{ width: '100%' }}
                placeholder="请输入退还租金（分）"
                min={0}
              />
            </Form.Item>
            <Form.Item
              name="refundDepositAmount"
              label="退还押金"
              rules={[
                { required: true, message: '请输入退还押金' },
                {
                  validator(_, value) {
                    if (!value || value <= (checkoutLease?.depositAmount || 0)) {
                      return Promise.resolve();
                    }
                    return Promise.reject(new Error('退还押金不能超过押金总额'));
                  },
                },
              ]}
            >
              <InputNumber
                style={{ width: '100%' }}
                placeholder="请输入退还押金（分）"
                min={0}
                max={checkoutLease?.depositAmount || 0}
              />
            </Form.Item>
          </div>

          {/* 收取费用 */}
          <div style={{ marginBottom: 16, padding: 16, border: '1px solid #d9d9d9', borderRadius: 8 }}>
            <h4 style={{ margin: '0 0 16px 0', color: '#52c41a' }}>收取费用</h4>
            <Form.Item
              name="waterAmount"
              label="水费（分）"
            >
              <InputNumber
                style={{ width: '100%' }}
                placeholder="请输入水费（分）"
                min={0}
              />
            </Form.Item>
            <Form.Item
              name="electricAmount"
              label="电费（分）"
            >
              <InputNumber
                style={{ width: '100%' }}
                placeholder="请输入电费（分）"
                min={0}
              />
            </Form.Item>
            <Form.Item
              name="otherAmount"
              label="其他费用（分）"
            >
              <InputNumber
                style={{ width: '100%' }}
                placeholder="请输入其他费用（分）"
                min={0}
              />
            </Form.Item>
          </div>

          {/* 备注 */}
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
        aggregateId={currentLease?.id || ''}
        domainType="lease"
        title={`租约操作日志 - ${currentLease?.tenantName || ''}`}
        onCancel={() => setOperationLogVisible(false)}
      />
    </div>
  );
};

export default Leases;
