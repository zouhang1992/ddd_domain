import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, message, Popconfirm, Pagination } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, HistoryOutlined, SearchOutlined, ReloadOutlined } from '@ant-design/icons';
import { locationApi, type LocationQueryParams, type LocationsQueryResult } from '../api/location';
import type { Location } from '../types/api';
import OperationLogModal from '../components/OperationLogModal';

const Locations: React.FC = () => {
  const [locations, setLocations] = useState<Location[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingLocation, setEditingLocation] = useState<Location | null>(null);
  const [form] = Form.useForm();
  const [operationLogVisible, setOperationLogVisible] = useState(false);
  const [currentLocation, setCurrentLocation] = useState<Location | null>(null);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [queryForm] = Form.useForm();

  const fetchLocations = async (params?: LocationQueryParams) => {
    setLoading(true);
    try {
      const queryParams = {
        ...params,
        offset: (page - 1) * pageSize,
        limit: pageSize,
      };
      const data: LocationsQueryResult = await locationApi.list(queryParams);
      setLocations(data.items);
      setTotal(data.total);
    } catch (error) {
      message.error('获取位置列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLocations();
  }, [page, pageSize]);

  const handleQuery = async () => {
    const values = await queryForm.validateFields();
    setPage(1);
    fetchLocations(values);
  };

  const handleReset = () => {
    queryForm.resetFields();
    setPage(1);
    fetchLocations();
  };

  const handlePageChange = (pageNum: number, pageSizeNum: number) => {
    setPage(pageNum);
    setPageSize(pageSizeNum);
  };

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

  const handleViewOperationLogs = (location: Location) => {
    setCurrentLocation(location);
    setOperationLogVisible(true);
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
        <h1>位置管理</h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新增位置
        </Button>
      </div>

      {/* 查询表单 */}
      <Form form={queryForm} layout="inline" style={{ marginBottom: 16 }}>
        <Form.Item name="shortName" label="简称">
          <Input placeholder="请输入简称" style={{ width: 150 }} />
        </Form.Item>
        <Form.Item name="detail" label="详情">
          <Input placeholder="请输入详情" style={{ width: 200 }} />
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
        dataSource={locations}
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

      <OperationLogModal
        visible={operationLogVisible}
        aggregateId={currentLocation?.id || ''}
        domainType="location"
        title={`位置操作日志 - ${currentLocation?.shortName || ''}`}
        onCancel={() => setOperationLogVisible(false)}
      />
    </div>
  );
};

export default Locations;
