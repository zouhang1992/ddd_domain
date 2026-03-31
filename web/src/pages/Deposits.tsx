import React, { useState, useEffect, useMemo } from 'react';
import { Table, Button, message, Select, Pagination, Tag, Form } from 'antd';
import { HistoryOutlined, SearchOutlined, ReloadOutlined } from '@ant-design/icons';
import { depositApi, type DepositQueryParams, type DepositsQueryResult } from '../api/deposit';
import { leaseApi } from '../api/lease';
import { roomApi } from '../api/room';
import { locationApi } from '../api/location';
import type { Deposit, Lease, Room, Location } from '../types/api';
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
  const [allDeposits, setAllDeposits] = useState<Deposit[]>([]);
  const [displayDeposits, setDisplayDeposits] = useState<Deposit[]>([]);
  const [leases, setLeases] = useState<Lease[]>([]);
  const [rooms, setRooms] = useState<Room[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);
  const [loading, setLoading] = useState(false);
  const [operationLogVisible, setOperationLogVisible] = useState(false);
  const [currentDeposit, setCurrentDeposit] = useState<Deposit | null>(null);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [queryForm] = Form.useForm();
  const [queryStatus, setQueryStatus] = useState<string>();
  const [queryLocationId, setQueryLocationId] = useState<string>();
  const [queryRoomId, setQueryRoomId] = useState<string>();

  const fetchDeposits = async () => {
    setLoading(true);
    try {
      const data: DepositsQueryResult = await depositApi.list({ limit: 1000 });
      setAllDeposits(data.items || []);
    } catch (error) {
      message.error('获取押金列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchLeases = async () => {
    try {
      const data = await leaseApi.list({ limit: 1000 });
      setLeases(data.items || []);
    } catch (error) {
      message.error('获取租约列表失败');
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

  useEffect(() => {
    fetchDeposits();
    fetchLeases();
    fetchRooms();
    fetchLocations();
  }, []);

  // 筛选用于查询的房间
  const queryRooms = useMemo(() => {
    if (queryLocationId) {
      return rooms.filter(room => room.locationId === queryLocationId);
    }
    return rooms;
  }, [rooms, queryLocationId]);

  // 应用筛选和分页
  useEffect(() => {
    let filtered = [...allDeposits];

    // 状态筛选
    if (queryStatus) {
      filtered = filtered.filter(d => d.status === queryStatus);
    }

    // 位置筛选
    if (queryLocationId) {
      const roomIds = rooms.filter(r => r.locationId === queryLocationId).map(r => r.id);
      const leaseIds = leases.filter(l => roomIds.includes(l.roomId)).map(l => l.id);
      filtered = filtered.filter(d => leaseIds.includes(d.leaseId));
    }

    // 房间筛选
    if (queryRoomId) {
      const leaseIds = leases.filter(l => l.roomId === queryRoomId).map(l => l.id);
      filtered = filtered.filter(d => leaseIds.includes(d.leaseId));
    }

    // 分页
    const start = (page - 1) * pageSize;
    const end = start + pageSize;
    setDisplayDeposits(filtered.slice(start, end));
    setTotal(filtered.length);
  }, [allDeposits, leases, rooms, page, pageSize, queryStatus, queryLocationId, queryRoomId]);

  const handleQuery = async () => {
    const values = await queryForm.validateFields();
    setQueryStatus(values.status);
    setPage(1);
  };

  const handleReset = () => {
    queryForm.resetFields();
    setQueryStatus(undefined);
    setQueryLocationId(undefined);
    setQueryRoomId(undefined);
    setPage(1);
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

  const getLocationInfo = (deposit: Deposit) => {
    const lease = leases.find(l => l.id === deposit.leaseId);
    if (!lease) return { location: '-', room: '-', address: '-' };
    const room = rooms.find(r => r.id === lease.roomId);
    if (!room) return { location: '-', room: '-', address: '-' };
    const location = locations.find(l => l.id === room.locationId);
    return {
      location: location?.shortName || '-',
      room: room.roomNumber,
      address: location?.detail || location?.shortName || '-',
    };
  };

  const columns = [
    {
      title: '位置',
      key: 'location',
      width: 120,
      render: (_: any, record: Deposit) => getLocationInfo(record).location,
    },
    {
      title: '房间',
      key: 'room',
      width: 120,
      render: (_: unknown, record: Deposit) => {
        const info = getLocationInfo(record);
        const lease = leases.find(l => l.id === record.leaseId);
        if (!lease) return info.room;
        const room = rooms.find(r => r.id === lease.roomId);
        if (!room) return info.room;
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
      render: (_: any, record: Deposit) => getLocationInfo(record).address,
    },
    {
      title: '租户',
      key: 'tenant',
      width: 100,
      render: (_: any, record: Deposit) => {
        const lease = leases.find(l => l.id === record.leaseId);
        return lease?.tenantName || '-';
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
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
      width: 100,
      render: formatAmount,
    },
    { title: '退还时间', dataIndex: 'refundedAt', key: 'refundedAt', width: 160 },
    { title: '扣款时间', dataIndex: 'deductedAt', key: 'deductedAt', width: 160 },
    { title: '备注', dataIndex: 'note', key: 'note', ellipsis: true },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', width: 160 },
    {
      title: '操作',
      key: 'actions',
      width: 120,
      fixed: 'right' as const,
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
        <Form.Item name="locationId" label="位置">
          <Select
            placeholder="请选择位置"
            style={{ width: 150 }}
            allowClear
            value={queryLocationId}
            onChange={(value) => {
              setQueryLocationId(value);
              setQueryRoomId(undefined);
              setPage(1);
            }}
          >
            {locations.map(location => (
              <Option key={location.id} value={location.id}>
                {location.shortName}
              </Option>
            ))}
          </Select>
        </Form.Item>
        <Form.Item name="roomId" label="房间">
          <Select
            placeholder="请选择房间"
            style={{ width: 150 }}
            allowClear
            value={queryRoomId}
            onChange={(value) => {
              setQueryRoomId(value);
              setPage(1);
            }}
          >
            {queryRooms.map(room => {
              const location = locations.find(l => l.id === room.locationId);
              return (
                <Option key={room.id} value={room.id}>
                  {location?.shortName} - {room.roomNumber}
                </Option>
              );
            })}
          </Select>
        </Form.Item>
        <Form.Item label="状态">
          <Select
            placeholder="请选择状态"
            style={{ width: 120 }}
            allowClear
            value={queryStatus}
            onChange={(value) => { setQueryStatus(value); setPage(1); }}
          >
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
        dataSource={displayDeposits}
        rowKey="id"
        loading={loading}
        pagination={false}
        scroll={{ x: 1800 }}
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
