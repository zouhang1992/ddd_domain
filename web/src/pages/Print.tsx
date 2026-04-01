import React, { useState, useEffect, useMemo } from 'react';
import { Card, Button, message, Tabs, Table, Tag, Spin, Space, Form, Input, Select, DatePicker, Pagination, Descriptions, Modal } from 'antd';
import { DownloadOutlined, HistoryOutlined, SearchOutlined, ReloadOutlined, EyeOutlined } from '@ant-design/icons';
import type { Lease, Bill, Room, Location, Landlord } from '../types/api';
import { printApi, type PrintJobQueryParams, type PrintJobsQueryResult, type PrintJob } from '../api/print';
import { leaseApi, type LeasesQueryResult } from '../api/lease';
import { billApi, type BillsQueryResult } from '../api/bill';
import { roomApi } from '../api/room';
import { locationApi } from '../api/location';
import { landlordApi } from '../api/landlord';
import OperationLogModal from '../components/OperationLogModal';

const { Option } = Select;

const Print: React.FC = () => {
  const [activeTab, setActiveTab] = useState<string>('bill');
  const [allBills, setAllBills] = useState<Bill[]>([]);
  const [allLeases, setAllLeases] = useState<Lease[]>([]);
  const [rooms, setRooms] = useState<Room[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);
  const [landlords, setLandlords] = useState<Landlord[]>([]);
  const [printJobs, setPrintJobs] = useState<PrintJob[]>([]);
  const [loading, setLoading] = useState(false);
  const [printing, setPrinting] = useState<string | null>(null);
  const [operationLogVisible, setOperationLogVisible] = useState(false);
  const [currentBill, setCurrentBill] = useState<Bill | null>(null);
  const [currentLease, setCurrentLease] = useState<Lease | null>(null);
  const [currentPrintJob, setCurrentPrintJob] = useState<PrintJob | null>(null);
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
  const [detailModalVisible, setDetailModalVisible] = useState(false);

  // 账单筛选状态
  const [billQueryType, setBillQueryType] = useState<string>();
  const [billQueryStatus, setBillQueryStatus] = useState<string>();
  const [billQueryLocationId, setBillQueryLocationId] = useState<string>();
  const [billQueryRoomId, setBillQueryRoomId] = useState<string>();

  // 租约筛选状态
  const [leaseQueryTenantName, setLeaseQueryTenantName] = useState<string>();
  const [leaseQueryTenantPhone, setLeaseQueryTenantPhone] = useState<string>();
  const [leaseQueryStatus, setLeaseQueryStatus] = useState<string>();
  const [leaseQueryLocationId, setLeaseQueryLocationId] = useState<string>();
  const [leaseQueryRoomId, setLeaseQueryRoomId] = useState<string>();

  const [displayBills, setDisplayBills] = useState<Bill[]>([]);
  const [displayLeases, setDisplayLeases] = useState<Lease[]>([]);

  const fetchAllBills = async () => {
    try {
      const data: BillsQueryResult = await billApi.list({ limit: 1000 });
      setAllBills(data.items || []);
    } catch (error) {
      message.error('获取账单列表失败');
    }
  };

  const fetchAllLeases = async () => {
    try {
      const data: LeasesQueryResult = await leaseApi.list({ limit: 1000 });
      setAllLeases(data.items || []);
    } catch (error) {
      console.error('获取所有租约失败', error);
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

  const fetchRooms = async () => {
    try {
      const data = await roomApi.list();
      setRooms(data.items || []);
    } catch (error) {
      message.error('获取房间列表失败');
    }
  };

  const fetchLocations = async () => {
    try {
      const data = await locationApi.list();
      setLocations(data.items || []);
    } catch (error) {
      message.error('获取位置列表失败');
    }
  };

  const fetchLandlords = async () => {
    try {
      const data = await landlordApi.list();
      setLandlords(data.items || []);
    } catch (error) {
      message.error('获取房东列表失败');
    }
  };

  useEffect(() => {
    fetchRooms();
    fetchLocations();
    fetchLandlords();
    fetchAllBills();
    fetchAllLeases();
  }, []);

  useEffect(() => {
    if (activeTab === 'jobs') {
      fetchPrintJobs();
    }
  }, [activeTab, printJobPage, pageSize]);

  // 账单查询用的房间列表
  const billQueryRooms = useMemo(() => {
    if (billQueryLocationId) {
      return rooms.filter(room => room.locationId === billQueryLocationId);
    }
    return rooms;
  }, [rooms, billQueryLocationId]);

  // 租约查询用的房间列表
  const leaseQueryRooms = useMemo(() => {
    if (leaseQueryLocationId) {
      return rooms.filter(room => room.locationId === leaseQueryLocationId);
    }
    return rooms;
  }, [rooms, leaseQueryLocationId]);

  // 应用账单筛选和分页
  useEffect(() => {
    let filtered = [...allBills];

    // 位置筛选
    if (billQueryLocationId) {
      const roomIds = rooms.filter(r => r.locationId === billQueryLocationId).map(r => r.id);
      const leaseIds = allLeases.filter(l => roomIds.includes(l.roomId)).map(l => l.id);
      filtered = filtered.filter(b => leaseIds.includes(b.leaseId));
    }

    // 房间筛选
    if (billQueryRoomId) {
      const leaseIds = allLeases.filter(l => l.roomId === billQueryRoomId).map(l => l.id);
      filtered = filtered.filter(b => leaseIds.includes(b.leaseId));
    }

    // 类型筛选
    if (billQueryType) {
      filtered = filtered.filter(b => b.type === billQueryType);
    }

    // 状态筛选
    if (billQueryStatus) {
      filtered = filtered.filter(b => b.status === billQueryStatus);
    }

    // 分页
    const start = (billPage - 1) * pageSize;
    const end = start + pageSize;
    setDisplayBills(filtered.slice(start, end));
    setBillTotal(filtered.length);
  }, [allBills, allLeases, rooms, billPage, pageSize, billQueryType, billQueryStatus, billQueryLocationId, billQueryRoomId]);

  // 应用租约筛选和分页
  useEffect(() => {
    let filtered = [...allLeases];

    // 租户姓名筛选
    if (leaseQueryTenantName) {
      filtered = filtered.filter(l => l.tenantName?.includes(leaseQueryTenantName));
    }

    // 租户电话筛选
    if (leaseQueryTenantPhone) {
      filtered = filtered.filter(l => l.tenantPhone?.includes(leaseQueryTenantPhone));
    }

    // 状态筛选
    if (leaseQueryStatus) {
      filtered = filtered.filter(l => l.status === leaseQueryStatus);
    }

    // 位置筛选
    if (leaseQueryLocationId) {
      const roomIds = rooms.filter(r => r.locationId === leaseQueryLocationId).map(r => r.id);
      filtered = filtered.filter(l => roomIds.includes(l.roomId));
    }

    // 房间筛选
    if (leaseQueryRoomId) {
      filtered = filtered.filter(l => l.roomId === leaseQueryRoomId);
    }

    // 分页
    const start = (leasePage - 1) * pageSize;
    const end = start + pageSize;
    setDisplayLeases(filtered.slice(start, end));
    setLeaseTotal(filtered.length);
  }, [allLeases, rooms, leasePage, pageSize, leaseQueryTenantName, leaseQueryTenantPhone, leaseQueryStatus, leaseQueryLocationId, leaseQueryRoomId]);

  const handleBillQuery = () => {
    setBillPage(1);
  };

  const handleBillReset = () => {
    billQueryForm.resetFields();
    setBillQueryType(undefined);
    setBillQueryStatus(undefined);
    setBillQueryLocationId(undefined);
    setBillQueryRoomId(undefined);
    setBillPage(1);
  };

  const handleLeaseQuery = () => {
    setLeasePage(1);
  };

  const handleLeaseReset = () => {
    leaseQueryForm.resetFields();
    setLeaseQueryTenantName(undefined);
    setLeaseQueryTenantPhone(undefined);
    setLeaseQueryStatus(undefined);
    setLeaseQueryLocationId(undefined);
    setLeaseQueryRoomId(undefined);
    setLeasePage(1);
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

  const handleViewPrintJobDetail = (job: PrintJob) => {
    setCurrentPrintJob(job);
    setDetailModalVisible(true);
  };

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

  const handleViewPrintJobOperationLogs = (job: PrintJob) => {
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

  const printJobTypeColorMap: Record<string, string> = {
    bill: 'blue',
    lease: 'green',
    invoice: 'purple',
  };

  const billColumns = [
    {
      title: '位置',
      key: 'location',
      width: 120,
      render: (_: any, record: Bill) => {
        const lease = allLeases.find(l => l.id === record.leaseId);
        if (!lease) return '-';
        const room = rooms.find(r => r.id === lease.roomId);
        const location = room ? locations.find(l => l.id === room.locationId) : null;
        return location?.shortName || '-';
      },
    },
    {
      title: '房间',
      key: 'room',
      width: 120,
      render: (_: unknown, record: Bill) => {
        const lease = allLeases.find(l => l.id === record.leaseId);
        if (!lease) return '-';
        const room = rooms.find(r => r.id === lease.roomId);
        if (!room) return '-';
        const location = locations.find(l => l.id === room.locationId);
        return (
          <span>
            <Tag color="blue" style={{ marginRight: 4, fontSize: '10px' }}>
              {location?.shortName || '未知位置'}
            </Tag>
            {room.roomNumber}
          </span>
        );
      },
    },
    {
      title: '地址',
      key: 'address',
      width: 150,
      ellipsis: true,
      render: (_: unknown, record: Bill) => {
        const lease = allLeases.find(l => l.id === record.leaseId);
        if (!lease) return '-';
        const room = rooms.find(r => r.id === lease.roomId);
        if (!room) return '-';
        const location = locations.find(l => l.id === room.locationId);
        return location?.detail || location?.shortName || '-';
      },
    },
    {
      title: '房东',
      key: 'landlord',
      width: 100,
      render: (_: unknown, record: Bill) => {
        const lease = allLeases.find(l => l.id === record.leaseId);
        if (!lease) return '-';
        return landlords.find(l => l.id === lease.landlordId)?.name || '-';
      },
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
      fixed: 'right' as const,
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
      title: '位置',
      key: 'location',
      width: 120,
      render: (_: any, record: Lease) => {
        const room = rooms.find(r => r.id === record.roomId);
        const location = room ? locations.find(l => l.id === room.locationId) : null;
        return location?.shortName || '-';
      },
    },
    {
      title: '房间',
      key: 'room',
      width: 120,
      render: (_: unknown, record: Lease) => {
        const room = rooms.find(r => r.id === record.roomId);
        if (!room) return '-';
        const location = locations.find(l => l.id === room.locationId);
        return (
          <span>
            <Tag color="blue" style={{ marginRight: 4, fontSize: '10px' }}>
              {location?.shortName || '未知位置'}
            </Tag>
            {room.roomNumber}
          </span>
        );
      },
    },
    {
      title: '地址',
      key: 'address',
      width: 150,
      ellipsis: true,
      render: (_: unknown, record: Lease) => {
        const room = rooms.find(r => r.id === record.roomId);
        if (!room) return '-';
        const location = locations.find(l => l.id === room.locationId);
        return location?.detail || location?.shortName || '-';
      },
    },
    {
      title: '房东',
      key: 'landlord',
      width: 100,
      render: (_: unknown, record: Lease) => {
        return landlords.find(l => l.id === record.landlordId)?.name || '-';
      },
    },
    {
      title: '租户姓名',
      dataIndex: 'tenantName',
      key: 'tenantName',
      width: 100,
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
      fixed: 'right' as const,
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
      title: '类型',
      dataIndex: 'type_text',
      key: 'type_text',
      width: 100,
      render: (text: string, record: PrintJob) => (
        <Tag color={printJobTypeColorMap[record.type] || 'default'}>
          {text}
        </Tag>
      ),
    },
    {
      title: '租户',
      dataIndex: 'tenant_name',
      key: 'tenant_name',
      width: 100,
    },
    {
      title: '房间号',
      dataIndex: 'room_number',
      key: 'room_number',
      width: 100,
      render: (text: string) => text || '-',
    },
    {
      title: '地址',
      dataIndex: 'address',
      key: 'address',
      width: 150,
      ellipsis: true,
      render: (text: string) => text || '-',
    },
    {
      title: '房东',
      dataIndex: 'landlord_name',
      key: 'landlord_name',
      width: 100,
      render: (text: string) => text || '-',
    },
    {
      title: '金额',
      dataIndex: 'amount_yuan',
      key: 'amount_yuan',
      width: 100,
      render: (amount: string) => amount ? `¥${amount}` : '-',
    },
    {
      title: '状态',
      dataIndex: 'status_text',
      key: 'status_text',
      width: 100,
      render: (text: string, record: PrintJob) => (
        <Tag color={printJobStatusColorMap[record.status] || 'default'}>
          {text}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 160,
    },
    {
      title: '操作',
      key: 'actions',
      width: 200,
      fixed: 'right' as const,
      render: (_: any, record: PrintJob) => (
        <Space size="small">
          <Button
            icon={<EyeOutlined />}
            size="small"
            onClick={() => handleViewPrintJobDetail(record)}
          >
            详情
          </Button>
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
            <Form.Item label="位置">
              <Select
                placeholder="请选择位置"
                style={{ width: 150 }}
                allowClear
                value={billQueryLocationId}
                onChange={(value) => {
                  setBillQueryLocationId(value);
                  setBillQueryRoomId(undefined);
                  setBillPage(1);
                }}
              >
                {locations.map(location => (
                  <Option key={location.id} value={location.id}>
                    {location.shortName}
                  </Option>
                ))}
              </Select>
            </Form.Item>
            <Form.Item label="房间">
              <Select
                placeholder="请选择房间"
                style={{ width: 150 }}
                allowClear
                value={billQueryRoomId}
                onChange={(value) => {
                  setBillQueryRoomId(value);
                  setBillPage(1);
                }}
              >
                {billQueryRooms.map(room => {
                  const location = locations.find(l => l.id === room.locationId);
                  return (
                    <Option key={room.id} value={room.id}>
                      {location?.shortName} - {room.roomNumber}
                    </Option>
                  );
                })}
              </Select>
            </Form.Item>
            <Form.Item label="类型">
              <Select
                placeholder="请选择类型"
                style={{ width: 120 }}
                allowClear
                value={billQueryType}
                onChange={(value) => { setBillQueryType(value); setBillPage(1); }}
              >
                <Option value="rent">租金</Option>
                <Option value="charge">收账</Option>
                <Option value="checkout">退租结算</Option>
              </Select>
            </Form.Item>
            <Form.Item label="状态">
              <Select
                placeholder="请选择状态"
                style={{ width: 120 }}
                allowClear
                value={billQueryStatus}
                onChange={(value) => { setBillQueryStatus(value); setBillPage(1); }}
              >
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
              dataSource={displayBills}
              columns={billColumns}
              rowKey="id"
              pagination={false}
              scroll={{ x: 1600 }}
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
            <Form.Item label="位置">
              <Select
                placeholder="请选择位置"
                style={{ width: 150 }}
                allowClear
                value={leaseQueryLocationId}
                onChange={(value) => {
                  setLeaseQueryLocationId(value);
                  setLeaseQueryRoomId(undefined);
                  setLeasePage(1);
                }}
              >
                {locations.map(location => (
                  <Option key={location.id} value={location.id}>
                    {location.shortName}
                  </Option>
                ))}
              </Select>
            </Form.Item>
            <Form.Item label="房间">
              <Select
                placeholder="请选择房间"
                style={{ width: 150 }}
                allowClear
                value={leaseQueryRoomId}
                onChange={(value) => {
                  setLeaseQueryRoomId(value);
                  setLeasePage(1);
                }}
              >
                {leaseQueryRooms.map(room => {
                  const location = locations.find(l => l.id === room.locationId);
                  return (
                    <Option key={room.id} value={room.id}>
                      {location?.shortName} - {room.roomNumber}
                    </Option>
                  );
                })}
              </Select>
            </Form.Item>
            <Form.Item label="租户姓名">
              <Input
                placeholder="请输入租户姓名"
                style={{ width: 120 }}
                value={leaseQueryTenantName}
                onChange={(e) => setLeaseQueryTenantName(e.target.value)}
              />
            </Form.Item>
            <Form.Item label="租户电话">
              <Input
                placeholder="请输入租户电话"
                style={{ width: 120 }}
                value={leaseQueryTenantPhone}
                onChange={(e) => setLeaseQueryTenantPhone(e.target.value)}
              />
            </Form.Item>
            <Form.Item label="状态">
              <Select
                placeholder="请选择状态"
                style={{ width: 120 }}
                allowClear
                value={leaseQueryStatus}
                onChange={(value) => { setLeaseQueryStatus(value); setLeasePage(1); }}
              >
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
              dataSource={displayLeases}
              columns={leaseColumns}
              rowKey="id"
              pagination={false}
              scroll={{ x: 1600 }}
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
              <Select placeholder="请选择状态" style={{ width: 120 }} allowClear>
                <Option value="pending">待处理</Option>
                <Option value="processing">处理中</Option>
                <Option value="completed">已完成</Option>
                <Option value="failed">失败</Option>
              </Select>
            </Form.Item>
            <Form.Item name="type" label="类型">
              <Select placeholder="请选择类型" style={{ width: 120 }} allowClear>
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
              scroll={{ x: 1400 }}
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

      {/* 打印作业详情弹窗 */}
      <Modal
        title="打印作业详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={[
          <Button key="close" onClick={() => setDetailModalVisible(false)}>
            关闭
          </Button>,
        ]}
        width={600}
      >
        {currentPrintJob && (
          <Descriptions bordered column={1}>
            <Descriptions.Item label="作业ID">
              {currentPrintJob.id}
            </Descriptions.Item>
            <Descriptions.Item label="类型">
              <Tag color={printJobTypeColorMap[currentPrintJob.type] || 'default'}>
                {currentPrintJob.type_text}
              </Tag>
            </Descriptions.Item>
            <Descriptions.Item label="状态">
              <Tag color={printJobStatusColorMap[currentPrintJob.status] || 'default'}>
                {currentPrintJob.status_text}
              </Tag>
            </Descriptions.Item>
            <Descriptions.Item label="租户姓名">
              {currentPrintJob.tenant_name}
            </Descriptions.Item>
            <Descriptions.Item label="租户电话">
              {currentPrintJob.tenant_phone || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="房间号">
              {currentPrintJob.room_number || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="地址">
              {currentPrintJob.address || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="房东姓名">
              {currentPrintJob.landlord_name || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="关联ID">
              {currentPrintJob.reference_id}
            </Descriptions.Item>
            <Descriptions.Item label="金额">
              {currentPrintJob.amount_yuan ? `¥${currentPrintJob.amount_yuan}` : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="创建时间">
              {currentPrintJob.created_at}
            </Descriptions.Item>
            <Descriptions.Item label="更新时间">
              {currentPrintJob.updated_at}
            </Descriptions.Item>
            {currentPrintJob.completed_at && (
              <Descriptions.Item label="完成时间">
                {currentPrintJob.completed_at}
              </Descriptions.Item>
            )}
            {currentPrintJob.error_msg && (
              <Descriptions.Item label="错误信息">
                <span style={{ color: '#d32f2f' }}>{currentPrintJob.error_msg}</span>
              </Descriptions.Item>
            )}
          </Descriptions>
        )}
      </Modal>

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
