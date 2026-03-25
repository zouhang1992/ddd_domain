import React, { useState, useEffect } from 'react';
import { Card, Button, message, Tabs, Table, Tag, Spin } from 'antd';
import { DownloadOutlined } from '@ant-design/icons';
import type { Lease, Bill } from '../types/api';
import { printApi } from '../api/print';
import { leaseApi } from '../api/lease';
import { billApi } from '../api/bill';

const Print: React.FC = () => {
  const [activeTab, setActiveTab] = useState<string>('bill');
  const [bills, setBills] = useState<Bill[]>([]);
  const [leases, setLeases] = useState<Lease[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [printing, setPrinting] = useState<string | null>(null);

  const fetchBills = async () => {
    setLoading(true);
    try {
      const data = await billApi.list();
      setBills(data || []);
    } catch (error) {
      message.error('获取账单列表失败');
    } finally {
      setLoading(false);
    }
  };

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

  useEffect(() => {
    if (activeTab === 'bill') {
      fetchBills();
    } else if (activeTab === 'lease') {
      fetchLeases();
    }
  }, [activeTab]);

  const handleDownloadReceipt = async (billId: string) => {
    setPrinting(billId);
    try {
      await printApi.downloadReceipt(billId);
      message.success('收据下载成功');
    } catch (error) {
      message.error('下载收据失败');
    } finally {
      setPrinting(null);
    }
  };

  const handleDownloadContract = async (leaseId: string) => {
    setPrinting(leaseId);
    try {
      await printApi.downloadContract(leaseId);
      message.success('合同下载成功');
    } catch (error) {
      message.error('下载合同失败');
    } finally {
      setPrinting(null);
    }
  };

  const statusColorMap: Record<string, string> = {
    paid: 'success',
    unpaid: 'default',
    pending: 'processing',
  };

  const billColumns = [
    {
      title: '账单ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      width: 100,
    },
    {
      title: '金额（元）',
      dataIndex: 'amount',
      key: 'amount',
      width: 120,
      render: (amount: number) => `¥${(amount / 100).toFixed(2)}`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => (
        <Tag color={statusColorMap[status] || 'default'}>
          {status === 'paid' ? '已支付' : status === 'unpaid' ? '未支付' : status}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      width: 150,
      render: (_: any, record: Bill) => (
        <Button
          type="primary"
          icon={<DownloadOutlined />}
          size="small"
          loading={printing === record.id}
          onClick={() => handleDownloadReceipt(record.id)}
        >
          下载收据
        </Button>
      ),
    },
  ];

  const leaseColumns = [
    {
      title: '租约ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '租户姓名',
      dataIndex: 'tenantName',
      key: 'tenantName',
    },
    {
      title: '开始日期',
      dataIndex: 'startDate',
      key: 'startDate',
      width: 120,
    },
    {
      title: '结束日期',
      dataIndex: 'endDate',
      key: 'endDate',
      width: 120,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => (
        <Tag color={status === 'active' ? 'success' : 'default'}>
          {status === 'active' ? '生效中' : status}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      width: 150,
      render: (_: any, record: Lease) => (
        <Button
          type="primary"
          icon={<DownloadOutlined />}
          size="small"
          loading={printing === record.id}
          onClick={() => handleDownloadContract(record.id)}
        >
          下载合同
        </Button>
      ),
    },
  ];

  const tabItems = [
    {
      key: 'bill',
      label: '打印账单',
      children: (
        <Card title="账单收据打印" style={{ marginTop: 16 }}>
          <Spin spinning={loading} tip="加载中...">
            <Table
              dataSource={bills}
              columns={billColumns}
              rowKey="id"
              pagination={{ pageSize: 10 }}
              scroll={{ x: 600 }}
            />
          </Spin>
        </Card>
      ),
    },
    {
      key: 'lease',
      label: '打印租约',
      children: (
        <Card title="租约合同打印" style={{ marginTop: 16 }}>
          <Spin spinning={loading} tip="加载中...">
            <Table
              dataSource={leases}
              columns={leaseColumns}
              rowKey="id"
              pagination={{ pageSize: 10 }}
              scroll={{ x: 600 }}
            />
          </Spin>
        </Card>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h1>打印服务</h1>
      </div>
      <Tabs activeKey={activeTab} onChange={setActiveTab} items={tabItems} />
    </div>
  );
};

export default Print;
