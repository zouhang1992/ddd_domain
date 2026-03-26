import React, { useState, useEffect } from 'react';
import { Modal, Table, Button, Space, Tag, message, Pagination } from 'antd';
import { ReloadOutlined } from '@ant-design/icons';
import { operationLogApi, type OperationLog } from '../api/operationLog';
import dayjs from 'dayjs';

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

interface OperationLogModalProps {
  visible: boolean;
  aggregateId: string;
  domainType: string;
  title?: string;
  onCancel: () => void;
}

const OperationLogModal: React.FC<OperationLogModalProps> = ({
  visible,
  aggregateId,
  domainType,
  title,
  onCancel,
}) => {
  const [logs, setLogs] = useState<OperationLog[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  const fetchLogs = async (pageNum: number = 1, pageSizeNum: number = 10) => {
    setLoading(true);
    try {
      const params = {
        aggregateId: aggregateId,
        domainType: domainType,
        offset: (pageNum - 1) * pageSizeNum,
        limit: pageSizeNum,
      };

      const result = await operationLogApi.list(params);
      setLogs(result.items || []);
      setTotal(result.total || 0);
      setPage(result.page || 1);
      setPageSize(result.limit || 10);
    } catch (error) {
      message.error('获取操作日志失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (visible && aggregateId && domainType) {
      fetchLogs(1, 10);
    }
  }, [visible, aggregateId, domainType]);

  const handlePageChange = (pageNum: number, pageSizeNum: number) => {
    fetchLogs(pageNum, pageSizeNum);
  };

  const handleRefresh = () => {
    fetchLogs(page, pageSize);
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
      title: '操作人',
      dataIndex: 'operatorId',
      key: 'operatorId',
      width: 120,
      render: (id?: string) => id || '-',
    },
  ];

  return (
    <Modal
      title={title || `操作日志 - ${domainType} ${aggregateId}`}
      open={visible}
      onCancel={onCancel}
      footer={[
        <Button key="close" onClick={onCancel}>
          关闭
        </Button>,
      ]}
      width={800}
    >
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <div>
          查看对象：{domainType} - {aggregateId}
        </div>
        <Space>
          <Button
            icon={<ReloadOutlined />}
            onClick={handleRefresh}
            loading={loading}
          >
            刷新
          </Button>
        </Space>
      </div>

      <Table
        columns={columns}
        dataSource={logs}
        rowKey="id"
        loading={loading}
        pagination={false}
        scroll={{ y: 400 }}
      />

      <div style={{ marginTop: 16 }}>
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
    </Modal>
  );
};

export default OperationLogModal;
