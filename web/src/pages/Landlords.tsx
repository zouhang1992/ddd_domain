import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, message, Popconfirm } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { landlordApi } from '../api/landlord';
import type { Landlord } from '../types/api';

const Landlords: React.FC = () => {
  const [landlords, setLandlords] = useState<Landlord[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingLandlord, setEditingLandlord] = useState<Landlord | null>(null);
  const [form] = Form.useForm();

  const fetchLandlords = async () => {
    setLoading(true);
    try {
      const data = await landlordApi.list();
      setLandlords(data);
    } catch (error) {
      message.error('获取房东列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLandlords();
  }, []);

  const handleCreate = () => {
    setEditingLandlord(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (landlord: Landlord) => {
    setEditingLandlord(landlord);
    form.setFieldsValue(landlord);
    setModalVisible(true);
  };

  const handleDelete = async (id: string) => {
    try {
      await landlordApi.delete(id);
      message.success('删除成功');
      fetchLandlords();
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      if (editingLandlord) {
        await landlordApi.update(editingLandlord.id, values);
        message.success('更新成功');
      } else {
        await landlordApi.create(values);
        message.success('创建成功');
      }
      setModalVisible(false);
      fetchLandlords();
    } catch (error) {
      message.error('操作失败');
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id' },
    { title: '姓名', dataIndex: 'name', key: 'name' },
    { title: '电话', dataIndex: 'phone', key: 'phone' },
    { title: '备注', dataIndex: 'note', key: 'note' },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: Landlord) => (
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
        <h1>房东管理</h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新增房东
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={landlords}
        rowKey="id"
        loading={loading}
      />

      <Modal
        title={editingLandlord ? '编辑房东' : '新增房东'}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="姓名"
            rules={[{ required: true, message: '请输入姓名' }]}
          >
            <Input placeholder="请输入姓名" />
          </Form.Item>
          <Form.Item
            name="phone"
            label="电话"
            rules={[{ required: true, message: '请输入电话' }]}
          >
            <Input placeholder="请输入电话" />
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

export default Landlords;
