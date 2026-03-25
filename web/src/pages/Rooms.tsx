import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Select, message, Popconfirm } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { roomApi } from '../api/room';
import { locationApi } from '../api/location';
import type { Room, Location } from '../types/api';

const { Option } = Select;

const Rooms: React.FC = () => {
  const [rooms, setRooms] = useState<Room[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingRoom, setEditingRoom] = useState<Room | null>(null);
  const [form] = Form.useForm();

  const fetchRooms = async () => {
    setLoading(true);
    try {
      const data = await roomApi.list();
      setRooms(data);
    } catch (error) {
      message.error('获取房间列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchLocations = async () => {
    try {
      const data = await locationApi.list();
      setLocations(data || []);
    } catch (error) {
      message.error('获取位置列表失败');
      setLocations([]);
    }
  };

  useEffect(() => {
    fetchRooms();
    fetchLocations();
  }, []);

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

      <Table
        columns={columns}
        dataSource={rooms}
        rowKey="id"
        loading={loading}
      />

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
    </div>
  );
};

export default Rooms;
