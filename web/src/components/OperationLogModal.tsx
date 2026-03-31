import React, { useState, useEffect } from 'react';
import { Modal, Table, Button, Space, Tag, message, Pagination } from 'antd';
import { ReloadOutlined, EyeOutlined } from '@ant-design/icons';
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
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [currentLog, setCurrentLog] = useState<OperationLog | null>(null);

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

  const handleViewDetail = (log: OperationLog) => {
    setCurrentLog(log);
    setDetailModalVisible(true);
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
    <>
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

      {/* 操作日志详情模态框 */}
      <Modal
        title="操作日志详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={[
          <Button key="close" onClick={() => setDetailModalVisible(false)}>
            关闭
          </Button>,
        ]}
        width={600}
      >
        {currentLog && (
          <div>
            <h3 style={{ marginBottom: 16, paddingBottom: 12, borderBottom: '1px solid #f0f0f0' }}>
              基本信息
            </h3>
            <table style={{ width: '100%', marginBottom: 20 }}>
              <tbody>
                <tr>
                  <td style={{ width: 120, padding: '8px 0', fontWeight: 'bold' }}>ID:</td>
                  <td style={{ padding: '8px 0' }}>{currentLog.id}</td>
                </tr>
                <tr>
                  <td style={{ padding: '8px 0', fontWeight: 'bold' }}>时间:</td>
                  <td style={{ padding: '8px 0' }}>{formatTimestamp(currentLog.timestamp)}</td>
                </tr>
                <tr>
                  <td style={{ padding: '8px 0', fontWeight: 'bold' }}>事件:</td>
                  <td style={{ padding: '8px 0' }}>
                    <Tag color="blue">{currentLog.eventName}</Tag>
                  </td>
                </tr>
                <tr>
                  <td style={{ padding: '8px 0', fontWeight: 'bold' }}>操作类型:</td>
                  <td style={{ padding: '8px 0' }}>
                    <Tag color={actionColorMap[currentLog.action] || 'default'}>
                      {currentLog.action}
                    </Tag>
                  </td>
                </tr>
                <tr>
                  <td style={{ padding: '8px 0', fontWeight: 'bold' }}>领域类型:</td>
                  <td style={{ padding: '8px 0' }}>{currentLog.domainType}</td>
                </tr>
                <tr>
                  <td style={{ padding: '8px 0', fontWeight: 'bold' }}>聚合ID:</td>
                  <td style={{ padding: '8px 0' }}>{currentLog.aggregateId || '-'}</td>
                </tr>
                <tr>
                  <td style={{ padding: '8px 0', fontWeight: 'bold' }}>操作人:</td>
                  <td style={{ padding: '8px 0' }}>{currentLog.operatorId || '-'}</td>
                </tr>
                <tr>
                  <td style={{ padding: '8px 0', fontWeight: 'bold' }}>创建时间:</td>
                  <td style={{ padding: '8px 0' }}>{formatTimestamp(currentLog.createdAt)}</td>
                </tr>
              </tbody>
            </table>

            {currentLog.details && Object.keys(currentLog.details).length > 0 && (
              <>
                <h3 style={{ marginBottom: 16, paddingBottom: 12, borderBottom: '1px solid #f0f0f0' }}>
                  详细数据
                </h3>
                <pre style={{
                  background: '#f5f5f5',
                  padding: 16,
                  borderRadius: 4,
                  overflow: 'auto',
                  maxHeight: 400,
                }}>
                  {JSON.stringify(currentLog.details, null, 2)}
                </pre>
              </>
            )}
          </div>
        )}
      </Modal>
    </>
  );
};

export default OperationLogModal;
