import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Select, DatePicker, message, Space, Tag } from 'antd';
import { SearchOutlined, ReloadOutlined, EyeOutlined } from '@ant-design/icons';
import { operationLogApi, type OperationLog } from '../api/operationLog';
import dayjs from 'dayjs';

const { Option } = Select;

// 领域类型选项
const domainTypeOptions = [
  { value: '', label: '全部' },
  { value: 'landlord', label: '房东' },
  { value: 'lease', label: '租约' },
  { value: 'bill', label: '账单' },
  { value: 'location', label: '位置' },
  { value: 'room', label: '房间' },
  { value: 'print', label: '打印' },
  { value: 'unknown', label: '未知' },
];

// 操作类型颜色映射
const actionColorMap: Record<string, string> = {
  created: 'green',
  updated: 'blue',
  deleted: 'red',
  renewed: 'orange',
  checkout: 'purple',
  activated: 'cyan',
  paid: 'lime',
  printed: 'gold',
  unknown: 'default',
};

// 领域类型颜色映射
const domainColorMap: Record<string, string> = {
  landlord: 'blue',
  lease: 'green',
  bill: 'orange',
  location: 'purple',
  room: 'cyan',
  print: 'gold',
  unknown: 'default',
};

const OperationLogs: React.FC = () => {
  const [logs, setLogs] = useState<OperationLog[]>([]);
  const [loading, setLoading] = useState(false);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [selectedLog, setSelectedLog] = useState<OperationLog | null>(null);
  const [form] = Form.useForm();
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);

  const fetchLogs = async (pageNum: number = 1, pageSizeNum: number = 20) => {
    setLoading(true);
    try {
      const values = form.getFieldsValue();
      const params: any = {
        offset: (pageNum - 1) * pageSizeNum,
        limit: pageSizeNum,
      };

      if (values.domainType) params.domainType = values.domainType;
      if (values.eventName) params.eventName = values.eventName;
      if (values.aggregateId) params.aggregateId = values.aggregateId;
      if (values.operatorId) params.operatorId = values.operatorId;
      if (values.timeRange && values.timeRange.length === 2) {
        params.startTime = values.timeRange[0].toISOString();
        params.endTime = values.timeRange[1].toISOString();
      }

      const result = await operationLogApi.list(params);
      setLogs(result.items || []);
      setTotal(result.total || 0);
      setPage(result.page || 1);
      setPageSize(result.limit || 20);
    } catch (error) {
      message.error('获取操作日志列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLogs(1, 20);
  }, []);

  const handleSearch = () => {
    fetchLogs(1, pageSize);
  };

  const handleReset = () => {
    form.resetFields();
    fetchLogs(1, 20);
  };

  const handleViewDetail = (log: OperationLog) => {
    setSelectedLog(log);
    setDetailModalVisible(true);
  };

  const handleTableChange = (pagination: any) => {
    fetchLogs(pagination.current, pagination.pageSize);
  };

  const formatTimestamp = (timestamp: string) => {
    return dayjs(timestamp).format('YYYY-MM-DD HH:mm:ss');
  };

  const columns = [
    {
      title: '时间',
      dataIndex: 'timestamp',
      key: 'timestamp',
      width: 180,
      render: formatTimestamp,
      sorter: (a: OperationLog, b: OperationLog) =>
        dayjs(a.timestamp).valueOf() - dayjs(b.timestamp).valueOf(),
    },
    {
      title: '事件名称',
      dataIndex: 'eventName',
      key: 'eventName',
      width: 180,
    },
    {
      title: '领域类型',
      dataIndex: 'domainType',
      key: 'domainType',
      width: 100,
      render: (domainType: string) => (
        <Tag color={domainColorMap[domainType] || 'default'}>
          {domainType}
        </Tag>
      ),
    },
    {
      title: '操作类型',
      dataIndex: 'action',
      key: 'action',
      width: 100,
      render: (action: string) => (
        <Tag color={actionColorMap[action] || 'default'}>
          {action}
        </Tag>
      ),
    },
    {
      title: '聚合ID',
      dataIndex: 'aggregateId',
      key: 'aggregateId',
      width: 120,
      render: (id?: string) => id || '-',
    },
    {
      title: '操作人',
      dataIndex: 'operatorId',
      key: 'operatorId',
      width: 120,
      render: (id?: string) => id || '-',
    },
    {
      title: '操作',
      key: 'actions',
      width: 100,
      render: (_: any, record: OperationLog) => (
        <Button
          icon={<EyeOutlined />}
          onClick={() => handleViewDetail(record)}
          size="small"
        >
          详情
        </Button>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <h1>操作日志</h1>
      </div>

      <Form form={form} layout="inline" style={{ marginBottom: 16 }}>
        <Form.Item name="domainType" label="领域类型">
          <Select placeholder="请选择" style={{ width: 120 }}>
            {domainTypeOptions.map(opt => (
              <Option key={opt.value} value={opt.value}>
                {opt.label}
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item name="eventName" label="事件名称">
          <Input placeholder="模糊搜索" style={{ width: 150 }} />
        </Form.Item>

        <Form.Item name="aggregateId" label="聚合ID">
          <Input placeholder="请输入" style={{ width: 120 }} />
        </Form.Item>

        <Form.Item name="timeRange" label="时间范围">
          <DatePicker.RangePicker
            showTime
            style={{ width: 350 }}
          />
        </Form.Item>

        <Form.Item>
          <Space>
            <Button
              type="primary"
              icon={<SearchOutlined />}
              onClick={handleSearch}
            >
              查询
            </Button>
            <Button
              icon={<ReloadOutlined />}
              onClick={handleReset}
            >
              重置
            </Button>
          </Space>
        </Form.Item>
      </Form>

      <Table
        columns={columns}
        dataSource={logs}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条`,
        }}
        onChange={handleTableChange}
        scroll={{ x: 1000 }}
      />

      <Modal
        title="操作日志详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={[
          <Button key="close" onClick={() => setDetailModalVisible(false)}>
            关闭
          </Button>,
        ]}
        width={700}
      >
        {selectedLog && (
          <div>
            <h3>基本信息</h3>
            <table style={{ width: '100%', marginBottom: 20 }}>
              <tbody>
                <tr>
                  <td style={{ width: 120, padding: '8px 0', fontWeight: 'bold' }}>ID:</td>
                  <td style={{ padding: '8px 0' }}>{selectedLog.id}</td>
                </tr>
                <tr>
                  <td style={{ padding: '8px 0', fontWeight: 'bold' }}>时间:</td>
                  <td style={{ padding: '8px 0' }}>{formatTimestamp(selectedLog.timestamp)}</td>
                </tr>
                <tr>
                  <td style={{ padding: '8px 0', fontWeight: 'bold' }}>事件:</td>
                  <td style={{ padding: '8px 0' }}>
                    <Tag color="blue">{selectedLog.eventName}</Tag>
                  </td>
                </tr>
                <tr>
                  <td style={{ padding: '8px 0', fontWeight: 'bold' }}>领域:</td>
                  <td style={{ padding: '8px 0' }}>
                    <Tag color={domainColorMap[selectedLog.domainType] || 'default'}>
                      {selectedLog.domainType}
                    </Tag>
                  </td>
                </tr>
                <tr>
                  <td style={{ padding: '8px 0', fontWeight: 'bold' }}>操作:</td>
                  <td style={{ padding: '8px 0' }}>
                    <Tag color={actionColorMap[selectedLog.action] || 'default'}>
                      {selectedLog.action}
                    </Tag>
                  </td>
                </tr>
                <tr>
                  <td style={{ padding: '8px 0', fontWeight: 'bold' }}>聚合ID:</td>
                  <td style={{ padding: '8px 0' }}>{selectedLog.aggregateId || '-'}</td>
                </tr>
                <tr>
                  <td style={{ padding: '8px 0', fontWeight: 'bold' }}>操作人:</td>
                  <td style={{ padding: '8px 0' }}>{selectedLog.operatorId || '-'}</td>
                </tr>
                <tr>
                  <td style={{ padding: '8px 0', fontWeight: 'bold' }}>创建时间:</td>
                  <td style={{ padding: '8px 0' }}>{formatTimestamp(selectedLog.createdAt)}</td>
                </tr>
              </tbody>
            </table>

            {selectedLog.details && Object.keys(selectedLog.details).length > 0 && (
              <>
                <h3>详细数据</h3>
                <pre style={{
                  background: '#f5f5f5',
                  padding: 16,
                  borderRadius: 4,
                  overflow: 'auto',
                  maxHeight: 400,
                }}>
                  {JSON.stringify(selectedLog.details, null, 2)}
                </pre>
              </>
            )}
          </div>
        )}
      </Modal>
    </div>
  );
};

export default OperationLogs;
