import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, message, Popconfirm, Pagination } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, HistoryOutlined, SearchOutlined, ReloadOutlined } from '@ant-design/icons';
import { landlordApi, type LandlordQueryParams, type LandlordsQueryResult } from '../api/landlord';
import type { Landlord } from '../types/api';
import OperationLogModal from '../components/OperationLogModal';

const Landlords: React.FC = () => {
  const [landlords, setLandlords] = useState<Landlord[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingLandlord, setEditingLandlord] = useState<Landlord | null>(null);
  const [form] = Form.useForm();
  const [operationLogVisible, setOperationLogVisible] = useState(false);
  const [currentLandlord, setCurrentLandlord] = useState<Landlord | null>(null);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [queryForm] = Form.useForm();

  const fetchLandlords = async (params?: LandlordQueryParams) => {
    setLoading(true);
    try {
      const queryParams = {
        ...params,
        offset: (page - 1) * pageSize,
        limit: pageSize,
      };
      const data: LandlordsQueryResult = await landlordApi.list(queryParams);
      setLandlords(data.items);
      setTotal(data.total);
    } catch (error) {
      message.error('获取房东列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLandlords();
  }, [page, pageSize]);

  const handleQuery = async () => {
    const values = await queryForm.validateFields();
    setPage(1);
    fetchLandlords(values);
  };

  const handleReset = () => {
    queryForm.resetFields();
    setPage(1);
    fetchLandlords();
  };

  const handlePageChange = (pageNum: number, pageSizeNum: number) => {
    setPage(pageNum);
    setPageSize(pageSizeNum);
  };

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

  const handleViewOperationLogs = (landlord: Landlord) => {
    setCurrentLandlord(landlord);
    setOperationLogVisible(true);
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
        <h1>房东管理</h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新增房东
        </Button>
      </div>

      {/* 查询表单 */}
      <Form form={queryForm} layout="inline" style={{ marginBottom: 16 }}>
        <Form.Item name="name" label="姓名">
          <Input placeholder="请输入姓名" style={{ width: 150 }} />
        </Form.Item>
        <Form.Item name="phone" label="电话">
          <Input placeholder="请输入电话" style={{ width: 150 }} />
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
        dataSource={landlords}
        rowKey="id"
        loading={loading}
        pagination={false}
        scroll={{ x: 600 }}
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

      <OperationLogModal
        visible={operationLogVisible}
        aggregateId={currentLandlord?.id || ''}
        domainType="landlord"
        title={`房东操作日志 - ${currentLandlord?.name || ''}`}
        onCancel={() => setOperationLogVisible(false)}
      />
    </div>
  );
};

export default Landlords;
