import React, { useState, useEffect } from 'react';
import { Table, Button, message, Select, Pagination, Tag, Form } from 'antd';
import { HistoryOutlined, SearchOutlined, ReloadOutlined } from '@ant-design/icons';
import { depositApi, type DepositQueryParams, type DepositsQueryResult } from '../api/deposit';
import { leaseApi } from '../api/lease';
import type { Deposit, Lease } from '../types/api';
import OperationLogModal from '../components/OperationLogModal';

const { Option } = Select;

const statusMap: Record<string, string> = {
  collected: '已收取',
  returning: '待退还',
  returned: '已退还',
};

const statusColorMap: Record<string, string> = {
  collected: 'green',
  returning: 'orange',
  returned: 'red',
};

const Deposits: React.FC = () => {
  const [deposits, setDeposits] = useState<Deposit[]>([]);
  const [leases, setLeases] = useState<Lease[]>([]);
  const [loading, setLoading] = useState(false);
  const [operationLogVisible, setOperationLogVisible] = useState(false);
  const [currentDeposit, setCurrentDeposit] = useState<Deposit | null>(null);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [queryForm] = Form.useForm();

  const fetchDeposits = async (params?: DepositQueryParams) => {
    setLoading(true);
    try {
      const queryParams = {
        ...params,
        offset: (page - 1) * pageSize,
        limit: pageSize,
      };
      const data: DepositsQueryResult = await depositApi.list(queryParams);
      setDeposits(data.items);
      setTotal(data.total);
    } catch (error) {
      message.error('获取押金列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchLeases = async () => {
    try {
      const data = await leaseApi.list();
      setLeases(data.items || []);
    } catch (error) {
      message.error('获取租约列表失败');
    }
  };

  useEffect(() => {
    fetchDeposits();
    fetchLeases();
  }, [page, pageSize]);

  const handleQuery = async () => {
    const values = await queryForm.validateFields();
    setPage(1);
    fetchDeposits(values);
  };

  const handleReset = () => {
    queryForm.resetFields();
    setPage(1);
    fetchDeposits();
  };

  const handlePageChange = (pageNum: number, pageSizeNum: number) => {
    setPage(pageNum);
    setPageSize(pageSizeNum);
  };

  const handleViewOperationLogs = (deposit: Deposit) => {
    setCurrentDeposit(deposit);
    setOperationLogVisible(true);
  };

  const formatAmount = (amount: number) => {
    return `¥${(amount / 100).toFixed(2)}`;
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    {
      title: '租约',
      dataIndex: 'leaseId',
      key: 'leaseId',
      render: (leaseId: string) => leases.find(l => l.id === leaseId)?.id || leaseId,
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
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
      render: formatAmount,
    },
    { title: '退还时间', dataIndex: 'refundedAt', key: 'refundedAt' },
    { title: '扣款时间', dataIndex: 'deductedAt', key: 'deductedAt' },
    { title: '备注', dataIndex: 'note', key: 'note' },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: Deposit) => (
        <Button
          icon={<HistoryOutlined />}
          onClick={() => handleViewOperationLogs(record)}
          size="small"
        >
          操作日志
        </Button>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h1>押金管理</h1>
      </div>

      {/* 查询表单 */}
      <Form form={queryForm} layout="inline" style={{ marginBottom: 16 }}>
        <Form.Item name="status" label="状态">
          <Select placeholder="请选择状态" style={{ width: 120 }}>
            <Option value="collected">已收取</Option>
            <Option value="returning">待退还</Option>
            <Option value="returned">已退还</Option>
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
        dataSource={deposits}
        rowKey="id"
        loading={loading}
        pagination={false}
        scroll={{ x: 1400 }}
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

      <OperationLogModal
        visible={operationLogVisible}
        aggregateId={currentDeposit?.id || ''}
        domainType="deposit"
        title={`押金操作日志 - ${currentDeposit?.id || ''}`}
        onCancel={() => setOperationLogVisible(false)}
      />
    </div>
  );
};

export default Deposits;
