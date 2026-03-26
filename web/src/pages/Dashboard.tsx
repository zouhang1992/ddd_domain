import React, { useState, useEffect } from 'react';
import { Card, Row, Col, Statistic, Spin, Alert } from 'antd';
import { HomeOutlined, EnvironmentOutlined, UserOutlined, FileTextOutlined, DollarOutlined } from '@ant-design/icons';
import { locationApi } from '../api/location';
import { roomApi } from '../api/room';
import { landlordApi } from '../api/landlord';
import { leaseApi } from '../api/lease';
import { billApi } from '../api/bill';
import { Link } from 'react-router-dom';

const Dashboard: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [stats, setStats] = useState({
    locations: 0,
    rooms: 0,
    landlords: 0,
    leases: 0,
    bills: 0,
  });

  const fetchStats = async () => {
    setLoading(true);
    setError(null);
    try {
      const [locations, rooms, landlords, leases, bills] = await Promise.all([
        locationApi.list(),
        roomApi.list(),
        landlordApi.list(),
        leaseApi.list(),
        billApi.list(),
      ]);
      setStats({
        locations: locations?.total || 0,
        rooms: rooms?.total || 0,
        landlords: landlords?.total || 0,
        leases: leases?.total || 0,
        bills: bills?.total || 0,
      });
    } catch (err) {
      setError('获取统计数据失败');
      console.error('Failed to fetch stats:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStats();
  }, []);

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" />
      </div>
    );
  }

  if (error) {
    return (
      <Alert
        message="错误"
        description={error}
        type="error"
        showIcon
        action={
          <a onClick={fetchStats}>重试</a>
        }
      />
    );
  }

  return (
    <div>
      <h1>仪表盘</h1>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card hoverable>
            <Link to="/locations">
              <Statistic
                title="位置"
                value={stats.locations}
                prefix={<EnvironmentOutlined />}
                valueStyle={{ color: '#3f8600' }}
              />
            </Link>
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card hoverable>
            <Link to="/rooms">
              <Statistic
                title="房间"
                value={stats.rooms}
                prefix={<HomeOutlined />}
                valueStyle={{ color: '#1890ff' }}
              />
            </Link>
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card hoverable>
            <Link to="/landlords">
              <Statistic
                title="房东"
                value={stats.landlords}
                prefix={<UserOutlined />}
                valueStyle={{ color: '#722ed1' }}
              />
            </Link>
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card hoverable>
            <Link to="/leases">
              <Statistic
                title="租约"
                value={stats.leases}
                prefix={<FileTextOutlined />}
                valueStyle={{ color: '#fa8c16' }}
              />
            </Link>
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card hoverable>
            <Link to="/bills">
              <Statistic
                title="账单"
                value={stats.bills}
                prefix={<DollarOutlined />}
                valueStyle={{ color: '#cf1322' }}
              />
            </Link>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard;
