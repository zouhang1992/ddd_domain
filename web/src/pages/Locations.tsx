import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, message, Popconfirm } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { locationApi } from '../api/location';
import type { Location } from '../types/api';

const Locations: React.FC = () => {
  const [locations, setLocations] = useState<Location[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingLocation, setEditingLocation] = useState<Location | null>(null);
  const [form] = Form.useForm();

  const fetchLocations = async () => {
    setLoading(true);
    try {
      const data = await locationApi.list();
      setLocations(data);
    } catch (error) {
      message.error('获取位置列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLocations();
  }, []);

  const handleCreate = () => {
    setEditingLocation(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (location: Location) => {
    setEditingLocation(location);
    form.setFieldsValue(location);
    setModalVisible(true);
  };

  const handleDelete = async (id: string) => {
    try {
      await locationApi.delete(id);
      message.success('删除成功');
      fetchLocations();
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      if (editingLocation) {
        await locationApi.update(editingLocation.id, values);
        message.success('更新成功');
      } else {
        await locationApi.create(values);
        message.success('创建成功');
      }
      setModalVisible(false);
      fetchLocations();
    } catch (error) {
      message.error('操作失败');
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id' },
    { title: '简称', dataIndex: 'shortName', key: 'shortName' },
    { title: '详情', dataIndex: 'detail', key: 'detail' },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: Location) => (
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
        <h1>位置管理</h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新增位置
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={locations}
        rowKey="id"
        loading={loading}
      />

      <Modal
        title={editingLocation ? '编辑位置' : '新增位置'}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="shortName"
            label="简称"
            rules={[{ required: true, message: '请输入简称' }]}
          >
            <Input placeholder="请输入简称" />
          </Form.Item>
          <Form.Item
            name="detail"
            label="详情"
            rules={[{ required: true, message: '请输入详情' }]}
          >
            <Input.TextArea placeholder="请输入详情" rows={3} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default Locations;
