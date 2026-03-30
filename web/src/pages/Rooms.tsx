import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Select, message, Popconfirm, Pagination, Tag } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, HistoryOutlined, SearchOutlined, ReloadOutlined } from '@ant-design/icons';
import { roomApi, type RoomQueryParams, type RoomsQueryResult } from '../api/room';
import { locationApi } from '../api/location';
import type { Room, Location } from '../types/api';
import OperationLogModal from '../components/OperationLogModal';

const { Option } = Select;

// 房间状态映射
const roomStatusMap: Record<string, { label: string; color: string }> = {
  available: { label: '可出租', color: 'green' },
  rented: { label: '已出租', color: 'red' },
  maintain: { label: '维护中', color: 'orange' },
};

// 获取状态显示
const getRoomStatusDisplay = (status: string) => {
  const config = roomStatusMap[status] || { label: status, color: 'default' };
  return <Tag color={config.color}>{config.label}</Tag>;
};

const Rooms: React.FC = () => {
  const [rooms, setRooms] = useState<Room[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingRoom, setEditingRoom] = useState<Room | null>(null);
  const [form] = Form.useForm();
  const [operationLogVisible, setOperationLogVisible] = useState(false);
  const [currentRoom, setCurrentRoom] = useState<Room | null>(null);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [queryForm] = Form.useForm();

  const fetchRooms = async (params?: RoomQueryParams) => {
    setLoading(true);
    try {
      const queryParams = {
        ...params,
        offset: (page - 1) * pageSize,
        limit: pageSize,
      };
      const data: RoomsQueryResult = await roomApi.list(queryParams);
      setRooms(data.items);
      setTotal(data.total);
    } catch (error) {
      message.error('获取房间列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchLocations = async () => {
    try {
      const data = await locationApi.list();
      setLocations(data.items || []);
    } catch (error) {
      message.error('获取位置列表失败');
      setLocations([]);
    }
  };

  useEffect(() => {
    fetchRooms();
    fetchLocations();
  }, [page, pageSize]);

  const handleQuery = async () => {
    const values = await queryForm.validateFields();
    setPage(1);
    fetchRooms(values);
  };

  const handleReset = () => {
    queryForm.resetFields();
    setPage(1);
    fetchRooms();
  };

  const handlePageChange = (pageNum: number, pageSizeNum: number) => {
    setPage(pageNum);
    setPageSize(pageSizeNum);
  };

  const handleCreate = () => {
    setEditingRoom(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (room: Room) => {
    setEditingRoom(room);
    form.setFieldsValue({
      locationId: room.locationId,
      roomNumber: room.roomNumber,
      tags: room.tags.join(','),
    });
    setModalVisible(true);
  };

  const handleDelete = async (id: string) => {
    try {
      await roomApi.delete(id);
      message.success('删除成功');
      fetchRooms();
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleViewOperationLogs = (room: Room) => {
    setCurrentRoom(room);
    setOperationLogVisible(true);
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      const tags = values.tags ? values.tags.split(',').map((tag: string) => tag.trim()) : [];
      if (editingRoom) {
        await roomApi.update(editingRoom.id, { ...values, tags });
        message.success('更新成功');
      } else {
        await roomApi.create({ ...values, tags });
        message.success('创建成功');
      }
      setModalVisible(false);
      fetchRooms();
    } catch (error) {
      message.error('操作失败');
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id' },
    {
      title: '位置',
      key: 'location',
      render: (_: any, record: Room) => {
        const location = locations.find(l => l.id === record.locationId);
        return location?.shortName || record.locationId;
      },
    },
    { title: '房间号', dataIndex: 'roomNumber', key: 'roomNumber' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => getRoomStatusDisplay(status),
    },
    { title: '标签', dataIndex: 'tags', key: 'tags', render: (tags: string[]) => tags.join(', ') },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: Room) => (
        <div>
          <Button icon={<EditOutlined />} onClick={() => handleEdit(record)} style={{ marginRight: 8 }}>
            编辑
          </Button>
          <Button
            icon={<HistoryOutlined />}
            onClick={() => handleViewOperationLogs(record)}
            style={{ marginRight: 8 }}
          >
            操作日志
          </Button>
          <Popconfirm
            title="确定删除?"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button icon={<DeleteOutlined />} danger>
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
        <h1>房间管理</h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新增房间
        </Button>
      </div>

      {/* 查询表单 */}
      <Form form={queryForm} layout="inline" style={{ marginBottom: 16 }}>
        <Form.Item name="locationId" label="位置">
          <Select placeholder="请选择位置" style={{ width: 150 }}>
            {locations.map(location => (
              <Option key={location.id} value={location.id}>
                {location.shortName}
              </Option>
            ))}
          </Select>
        </Form.Item>
        <Form.Item name="roomNumber" label="房间号">
          <Input placeholder="请输入房间号" style={{ width: 120 }} />
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
        dataSource={rooms}
        rowKey="id"
        loading={loading}
        pagination={false}
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
        title={editingRoom ? '编辑房间' : '新增房间'}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="locationId"
            label="位置"
            rules={[{ required: true, message: '请选择位置' }]}
          >
            <Select placeholder="请选择位置">
              {(locations || []).map((location) => (
                <Option key={location.id} value={location.id}>
                  {location.shortName}
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="roomNumber"
            label="房间号"
            rules={[{ required: true, message: '请输入房间号' }]}
          >
            <Input placeholder="请输入房间号" />
          </Form.Item>
          <Form.Item
            name="tags"
            label="标签"
            help="多个标签用逗号分隔"
          >
            <Input placeholder="请输入标签，多个标签用逗号分隔" />
          </Form.Item>
        </Form>
      </Modal>

      <OperationLogModal
        visible={operationLogVisible}
        aggregateId={currentRoom?.id || ''}
        domainType="room"
        title={`房间操作日志 - ${currentRoom?.roomNumber || ''}`}
        onCancel={() => setOperationLogVisible(false)}
      />
    </div>
  );
};

export default Rooms;
