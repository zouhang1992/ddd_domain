import React, { useState, useEffect } from 'react';
import { Card, Button, message, Tabs, Table, Tag, Spin, Space, Form, Input, Select, DatePicker, Pagination } from 'antd';
import { DownloadOutlined, HistoryOutlined, SearchOutlined, ReloadOutlined } from '@ant-design/icons';
import type { Lease, Bill } from '../types/api';
import { printApi, type PrintJobQueryParams, type PrintJobsQueryResult } from '../api/print';
import { leaseApi, type LeaseQueryParams, type LeasesQueryResult } from '../api/lease';
import { billApi, type BillQueryParams, type BillsQueryResult } from '../api/bill';
import OperationLogModal from '../components/OperationLogModal';

const { Option } = Select;

const Print: React.FC = () => {
  const [activeTab, setActiveTab] = useState<string>('bill');
  const [bills, setBills] = useState<Bill[]>([]);
  const [leases, setLeases] = useState<Lease[]>([]);
  const [printJobs, setPrintJobs] = useState<any[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [printing, setPrinting] = useState<string | null>(null);
  const [operationLogVisible, setOperationLogVisible] = useState(false);
  const [currentBill, setCurrentBill] = useState<Bill | null>(null);
  const [currentLease, setCurrentLease] = useState<Lease | null>(null);
  const [currentPrintJob, setCurrentPrintJob] = useState<any | null>(null);
  const [billTotal, setBillTotal] = useState<number>(0);
  const [leaseTotal, setLeaseTotal] = useState<number>(0);
  const [printJobTotal, setPrintJobTotal] = useState<number>(0);
  const [billPage, setBillPage] = useState<number>(1);
  const [leasePage, setLeasePage] = useState<number>(1);
  const [printJobPage, setPrintJobPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);
  const [billQueryForm] = Form.useForm();
  const [leaseQueryForm] = Form.useForm();
  const [printJobQueryForm] = Form.useForm();

  const fetchBills = async (params?: BillQueryParams) => {
    setLoading(true);
    try {
      const queryParams = {
        ...params,
        offset: (billPage - 1) * pageSize,
        limit: pageSize,
      };
      const data: BillsQueryResult = await billApi.list(queryParams);
      setBills(data.items || []);
      setBillTotal(data.total);
    } catch (error) {
      message.error('获取账单列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchLeases = async (params?: LeaseQueryParams) => {
    setLoading(true);
    try {
      const queryParams = {
        ...params,
        offset: (leasePage - 1) * pageSize,
        limit: pageSize,
      };
      const data: LeasesQueryResult = await leaseApi.list(queryParams);
      setLeases(data.items || []);
      setLeaseTotal(data.total);
    } catch (error) {
      message.error('获取租约列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchPrintJobs = async (params?: PrintJobQueryParams) => {
    setLoading(true);
    try {
      const queryParams = {
        ...params,
        offset: (printJobPage - 1) * pageSize,
        limit: pageSize,
      };
      const data: PrintJobsQueryResult = await printApi.listPrintJobs(queryParams);
      setPrintJobs(data.items || []);
      setPrintJobTotal(data.total);
    } catch (error) {
      message.error('获取打印作业列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleBillQuery = async () => {
    const values = await billQueryForm.validateFields();
    setBillPage(1);
    fetchBills(values);
  };

  const handleBillReset = () => {
    billQueryForm.resetFields();
    setBillPage(1);
    fetchBills();
  };

  const handleLeaseQuery = async () => {
    const values = await leaseQueryForm.validateFields();
    setLeasePage(1);
    fetchLeases(values);
  };

  const handleLeaseReset = () => {
    leaseQueryForm.resetFields();
    setLeasePage(1);
    fetchLeases();
  };

  const handlePrintJobQuery = async () => {
    const values = await printJobQueryForm.validateFields();
    setPrintJobPage(1);
    fetchPrintJobs(values);
  };

  const handlePrintJobReset = () => {
    printJobQueryForm.resetFields();
    setPrintJobPage(1);
    fetchPrintJobs();
  };

  const handleBillPageChange = (pageNum: number, pageSizeNum: number) => {
    setBillPage(pageNum);
    setPageSize(pageSizeNum);
  };

  const handleLeasePageChange = (pageNum: number, pageSizeNum: number) => {
    setLeasePage(pageNum);
    setPageSize(pageSizeNum);
  };

  const handlePrintJobPageChange = (pageNum: number, pageSizeNum: number) => {
    setPrintJobPage(pageNum);
    setPageSize(pageSizeNum);
  };

  useEffect(() => {
    if (activeTab === 'bill') {
      fetchBills();
    } else if (activeTab === 'lease') {
      fetchLeases();
    } else if (activeTab === 'jobs') {
      fetchPrintJobs();
    }
  }, [activeTab, billPage, leasePage, printJobPage, pageSize]);

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

  const handleViewBillOperationLogs = (bill: Bill) => {
    setCurrentBill(bill);
    setCurrentLease(null);
    setCurrentPrintJob(null);
    setOperationLogVisible(true);
  };

  const handleViewLeaseOperationLogs = (lease: Lease) => {
    setCurrentLease(lease);
    setCurrentBill(null);
    setCurrentPrintJob(null);
    setOperationLogVisible(true);
  };

  const handleViewPrintJobOperationLogs = (job: any) => {
    setCurrentPrintJob(job);
    setCurrentBill(null);
    setCurrentLease(null);
    setOperationLogVisible(true);
  };

  const billStatusColorMap: Record<string, string> = {
    paid: 'success',
    unpaid: 'default',
    pending: 'processing',
  };

  const billStatusMap: Record<string, string> = {
    paid: '已支付',
    unpaid: '未支付',
    pending: '待支付',
  };

  const leaseStatusColorMap: Record<string, string> = {
    active: 'success',
    pending: 'default',
    expired: 'warning',
    checkout: 'error',
  };

  const leaseStatusMap: Record<string, string> = {
    active: '生效中',
    pending: '待生效',
    expired: '已过期',
    checkout: '已退租',
  };

  const printJobStatusColorMap: Record<string, string> = {
    pending: 'default',
    processing: 'processing',
    completed: 'success',
    failed: 'error',
  };

  const printJobStatusMap: Record<string, string> = {
    pending: '待处理',
    processing: '处理中',
    completed: '已完成',
    failed: '失败',
  };

  const printJobTypeColorMap: Record<string, string> = {
    bill: 'blue',
    lease: 'green',
    invoice: 'purple',
  };

  const printJobTypeMap: Record<string, string> = {
    bill: '账单收据',
    lease: '租约合同',
    invoice: '发票',
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
      render: (type: string) => {
        const typeMap: Record<string, string> = {
          charge: '收账',
          checkout: '退租结算',
        };
        return typeMap[type] || type;
      },
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
        <Tag color={billStatusColorMap[status] || 'default'}>
          {billStatusMap[status] || status}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      width: 200,
      render: (_: any, record: Bill) => (
        <Space size="small">
          <Button
            type="primary"
            icon={<DownloadOutlined />}
            size="small"
            loading={printing === record.id}
            onClick={() => handleDownloadReceipt(record.id)}
          >
            下载收据
          </Button>
          <Button
            icon={<HistoryOutlined />}
            size="small"
            onClick={() => handleViewBillOperationLogs(record)}
          >
            操作日志
          </Button>
        </Space>
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
        <Tag color={leaseStatusColorMap[status] || 'default'}>
          {leaseStatusMap[status] || status}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      width: 200,
      render: (_: any, record: Lease) => (
        <Space size="small">
          <Button
            type="primary"
            icon={<DownloadOutlined />}
            size="small"
            loading={printing === record.id}
            onClick={() => handleDownloadContract(record.id)}
          >
            下载合同
          </Button>
          <Button
            icon={<HistoryOutlined />}
            size="small"
            onClick={() => handleViewLeaseOperationLogs(record)}
          >
            操作日志
          </Button>
        </Space>
      ),
    },
  ];

  const printJobColumns = [
    {
      title: '作业ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      width: 100,
      render: (type: string) => (
        <Tag color={printJobTypeColorMap[type] || 'default'}>
          {printJobTypeMap[type] || type}
        </Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => (
        <Tag color={printJobStatusColorMap[status] || 'default'}>
          {printJobStatusMap[status] || status}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 160,
    },
    {
      title: '操作',
      key: 'actions',
      width: 150,
      render: (_: any, record: any) => (
        <Space size="small">
          <Button
            icon={<HistoryOutlined />}
            size="small"
            onClick={() => handleViewPrintJobOperationLogs(record)}
          >
            操作日志
          </Button>
        </Space>
      ),
    },
  ];

  const tabItems = [
    {
      key: 'bill',
      label: '打印账单',
      children: (
        <Card title="账单收据打印" style={{ marginTop: 16 }}>
          {/* 查询表单 */}
          <Form form={billQueryForm} layout="inline" style={{ marginBottom: 16 }}>
            <Form.Item name="type" label="类型">
              <Select placeholder="请选择类型" style={{ width: 120 }}>
                <Option value="charge">收账</Option>
                <Option value="checkout">退租结算</Option>
              </Select>
            </Form.Item>
            <Form.Item name="status" label="状态">
              <Select placeholder="请选择状态" style={{ width: 120 }}>
                <Option value="pending">待支付</Option>
                <Option value="paid">已支付</Option>
              </Select>
            </Form.Item>
            <Form.Item>
              <Button type="primary" icon={<SearchOutlined />} onClick={handleBillQuery} loading={loading}>
                查询
              </Button>
            </Form.Item>
            <Form.Item>
              <Button icon={<ReloadOutlined />} onClick={handleBillReset} loading={loading}>
                重置
              </Button>
            </Form.Item>
          </Form>
          <Spin spinning={loading} tip="加载中...">
            <Table
              dataSource={bills}
              columns={billColumns}
              rowKey="id"
              pagination={false}
              scroll={{ x: 600 }}
            />
          </Spin>
          <div style={{ marginTop: 16, display: 'flex', justifyContent: 'flex-end' }}>
            <Pagination
              current={billPage}
              pageSize={pageSize}
              total={billTotal}
              onChange={handleBillPageChange}
              showSizeChanger
              showQuickJumper
              showTotal={(total) => `共 ${total} 条记录`}
            />
          </div>
        </Card>
      ),
    },
    {
      key: 'lease',
      label: '打印租约',
      children: (
        <Card title="租约合同打印" style={{ marginTop: 16 }}>
          {/* 查询表单 */}
          <Form form={leaseQueryForm} layout="inline" style={{ marginBottom: 16 }}>
            <Form.Item name="tenantName" label="租户姓名">
              <Input placeholder="请输入租户姓名" style={{ width: 120 }} />
            </Form.Item>
            <Form.Item name="tenantPhone" label="租户电话">
              <Input placeholder="请输入租户电话" style={{ width: 120 }} />
            </Form.Item>
            <Form.Item name="status" label="状态">
              <Select placeholder="请选择状态" style={{ width: 120 }}>
                <Option value="pending">待生效</Option>
                <Option value="active">生效中</Option>
                <Option value="expired">已过期</Option>
                <Option value="checkout">已退租</Option>
              </Select>
            </Form.Item>
            <Form.Item>
              <Button type="primary" icon={<SearchOutlined />} onClick={handleLeaseQuery} loading={loading}>
                查询
              </Button>
            </Form.Item>
            <Form.Item>
              <Button icon={<ReloadOutlined />} onClick={handleLeaseReset} loading={loading}>
                重置
              </Button>
            </Form.Item>
          </Form>
          <Spin spinning={loading} tip="加载中...">
            <Table
              dataSource={leases}
              columns={leaseColumns}
              rowKey="id"
              pagination={false}
              scroll={{ x: 600 }}
            />
          </Spin>
          <div style={{ marginTop: 16, display: 'flex', justifyContent: 'flex-end' }}>
            <Pagination
              current={leasePage}
              pageSize={pageSize}
              total={leaseTotal}
              onChange={handleLeasePageChange}
              showSizeChanger
              showQuickJumper
              showTotal={(total) => `共 ${total} 条记录`}
            />
          </div>
        </Card>
      ),
    },
    {
      key: 'jobs',
      label: '打印作业',
      children: (
        <Card title="打印作业管理" style={{ marginTop: 16 }}>
          {/* 查询表单 */}
          <Form form={printJobQueryForm} layout="inline" style={{ marginBottom: 16 }}>
            <Form.Item name="status" label="状态">
              <Select placeholder="请选择状态" style={{ width: 120 }}>
                <Option value="pending">待处理</Option>
                <Option value="processing">处理中</Option>
                <Option value="completed">已完成</Option>
                <Option value="failed">失败</Option>
              </Select>
            </Form.Item>
            <Form.Item name="type" label="类型">
              <Select placeholder="请选择类型" style={{ width: 120 }}>
                <Option value="bill">账单收据</Option>
                <Option value="lease">租约合同</Option>
                <Option value="invoice">发票</Option>
              </Select>
            </Form.Item>
            <Form.Item name="startDate" label="开始日期">
              <DatePicker style={{ width: 150 }} placeholder="请选择开始日期" />
            </Form.Item>
            <Form.Item name="endDate" label="结束日期">
              <DatePicker style={{ width: 150 }} placeholder="请选择结束日期" />
            </Form.Item>
            <Form.Item>
              <Button type="primary" icon={<SearchOutlined />} onClick={handlePrintJobQuery} loading={loading}>
                查询
              </Button>
            </Form.Item>
            <Form.Item>
              <Button icon={<ReloadOutlined />} onClick={handlePrintJobReset} loading={loading}>
                重置
              </Button>
            </Form.Item>
          </Form>
          <Spin spinning={loading} tip="加载中...">
            <Table
              dataSource={printJobs}
              columns={printJobColumns}
              rowKey="id"
              pagination={false}
              scroll={{ x: 600 }}
            />
          </Spin>
          <div style={{ marginTop: 16, display: 'flex', justifyContent: 'flex-end' }}>
            <Pagination
              current={printJobPage}
              pageSize={pageSize}
              total={printJobTotal}
              onChange={handlePrintJobPageChange}
              showSizeChanger
              showQuickJumper
              showTotal={(total) => `共 ${total} 条记录`}
            />
          </div>
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

      <OperationLogModal
        visible={operationLogVisible}
        aggregateId={currentBill?.id || currentLease?.id || currentPrintJob?.id || ''}
        domainType={currentBill ? 'bill' : currentLease ? 'lease' : 'print'}
        title={
          currentBill
            ? `账单操作日志 - ${currentBill.id}`
            : currentLease
            ? `租约操作日志 - ${currentLease.tenantName}`
            : currentPrintJob
            ? `打印作业操作日志 - ${currentPrintJob.id}`
            : '操作日志'
        }
        onCancel={() => setOperationLogVisible(false)}
      />
    </div>
  );
};

export default Print;
