import React, { useState, useEffect, useMemo } from 'react';
import { Table, Button, Modal, Form, Input, message, Popconfirm, Select, DatePicker, InputNumber, Tag, Space } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, CheckOutlined, ReloadOutlined, PrinterOutlined, ThunderboltOutlined } from '@ant-design/icons';
import type { Lease, Room, Landlord, Location } from '../types/api';
import { leaseApi } from '../api/lease';
import { roomApi } from '../api/room';
import { landlordApi } from '../api/landlord';
import { locationApi } from '../api/location';
import dayjs from 'dayjs';

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
  const [leases, setLeases] = useState<Lease[]>([]);
  const [rooms, setRooms] = useState<Room[]>([]);
  const [landlords, setLandlords] = useState<Landlord[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [renewModalVisible, setRenewModalVisible] = useState(false);
  const [editingLease, setEditingLease] = useState<Lease | null>(null);
  const [form] = Form.useForm();
  const [renewForm] = Form.useForm();
  const [selectedLocationId, setSelectedLocationId] = useState<string>('');

  const fetchLeases = async () => {
    setLoading(true);
    try {
      const data = await leaseApi.list();
      setLeases(data || []);
    } catch (error) {
      message.error('获取租约列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchRooms = async () => {
    try {
      const data = await roomApi.list();
      setRooms(data || []);
    } catch (error) {
      message.error('获取房间列表失败');
    }
  };

  const fetchLandlords = async () => {
    try {
      const data = await landlordApi.list();
      setLandlords(data || []);
    } catch (error) {
      message.error('获取房东列表失败');
    }
  };

  const fetchLocations = async () => {
    try {
      const data = await locationApi.list();
      setLocations(data || []);
    } catch (error) {
      message.error('获取位置列表失败');
    }
  };

  useEffect(() => {
    fetchLeases();
    fetchRooms();
    fetchLandlords();
    fetchLocations();
  }, []);

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

  const handleCheckout = async (lease: Lease) => {
    try {
      await leaseApi.checkout(lease.id);
      message.success('退租成功');
      fetchLeases();
    } catch (error) {
      message.error('退租失败');
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await leaseApi.delete(id);
      message.success('删除成功');
      fetchLeases();
    } catch (error) {
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
    } catch (error) {
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
    } catch (error) {
      message.error('续租失败');
    }
  };

  const formatAmount = (amount: number) => {
    return `¥${(amount / 100).toFixed(2)}`;
  };

  // 根据选择的位置筛选房间
  const filteredRooms = useMemo(() => {
    if (!selectedLocationId) {
      return rooms;
    }
    return rooms.filter(room => room.locationId === selectedLocationId);
  }, [rooms, selectedLocationId]);

  // 处理位置选择变化
  const handleLocationChange = (value: string) => {
    setSelectedLocationId(value);
    form.setFieldValue('roomId', undefined);
  };

  const handlePrintContract = async (lease: Lease) => {
    try {
      await leaseApi.printContract(lease.id);
      message.success('合同下载成功');
    } catch (error) {
      message.error('合同下载失败');
    }
  };

  const handleActivateLease = async (lease: Lease) => {
    try {
      await leaseApi.activate(lease.id);
      message.success('租约生效成功');
      fetchLeases();
    } catch (error) {
      message.error('租约生效失败');
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
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
      dataIndex: 'landlordId',
      key: 'landlordId',
      render: (landlordId: string) => landlords.find(l => l.id === landlordId)?.name || landlordId,
    },
    { title: '租户姓名', dataIndex: 'tenantName', key: 'tenantName' },
    { title: '租户电话', dataIndex: 'tenantPhone', key: 'tenantPhone' },
    { title: '开始日期', dataIndex: 'startDate', key: 'startDate' },
    { title: '结束日期', dataIndex: 'endDate', key: 'endDate' },
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

      <Table
        columns={columns}
        dataSource={leases}
        rowKey="id"
        loading={loading}
        scroll={{ x: 1800 }}
      />

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
                <Select placeholder="请先选择位置，再选择房间">
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
    </div>
  );
};

export default Leases;
