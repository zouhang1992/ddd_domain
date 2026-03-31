import React, { useState, useEffect } from 'react';
import { Card, DatePicker, Statistic, Table, Button, Spin, Alert, message, Row, Col, Tag, Select } from 'antd';
import { DollarOutlined, DownloadOutlined, SearchOutlined } from '@ant-design/icons';
import { incomeApi } from '../api/income';
import { locationApi } from '../api/location';
import type { Location } from '../types/api';

const { MonthPicker } = DatePicker;
const { Option } = Select;

interface IncomeItem {
  key: string;
  type: string;
  amount: number;
  formattedAmount: string;
  percentage: string;
  isExpense?: boolean;
}

// 安全计算占比函数
const calculatePercentage = (value: number, total: number): string => {
  if (!total || total === 0 || isNaN(total)) {
    return '0%';
  }
  const absoluteValue = Math.abs(value);
  const percentage = (absoluteValue / total) * 100;
  if (isNaN(percentage) || !isFinite(percentage)) {
    return '0%';
  }
  return percentage.toFixed(1) + '%';
};

const Income: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [incomeData, setIncomeData] = useState<any>(null);
  const [selectedMonth, setSelectedMonth] = useState<string | null>(null);
  const [selectedLocation, setSelectedLocation] = useState<string | null>(null);
  const [locations, setLocations] = useState<Location[]>([]);

  const fetchIncome = async (month?: string, locationId?: string) => {
    setLoading(true);
    setError(null);
    try {
      const data = await incomeApi.getReport(month, locationId || undefined);
      setIncomeData(data);
      message.success('获取收入数据成功');
    } catch (err) {
      setError('获取收入数据失败');
      console.error('Failed to fetch income:', err);
    } finally {
      setLoading(false);
    }
  };

  const fetchLocations = async () => {
    try {
      const data = await locationApi.list();
      setLocations(data.items || []);
    } catch (err) {
      console.error('Failed to fetch locations:', err);
    }
  };

  useEffect(() => {
    fetchIncome();
    fetchLocations();
  }, []);

  const handleMonthChange = (date: any) => {
    if (date) {
      const monthStr = date.format('YYYY-MM');
      setSelectedMonth(monthStr);
      fetchIncome(monthStr, selectedLocation || undefined);
    } else {
      setSelectedMonth(null);
      fetchIncome(undefined, selectedLocation || undefined);
    }
  };

  const handleLocationChange = (locationId: string | null) => {
    setSelectedLocation(locationId);
    fetchIncome(selectedMonth || undefined, locationId || undefined);
  };

  const handleExport = () => {
    message.info('导出功能开发中...');
  };

  // 安全获取数值，避免 NaN
  const safeNumber = (value: any, defaultValue: number = 0): number => {
    const num = Number(value);
    return isNaN(num) ? defaultValue : num;
  };

  // 计算用于占比计算的有效总金额（取绝对值）
  const getEffectiveTotal = (data: any): number => {
    if (!data) return 0;
    const rentIncome = safeNumber(data.rentIncome);
    const waterIncome = safeNumber(data.waterIncome);
    const electricIncome = safeNumber(data.electricIncome);
    const otherIncome = safeNumber(data.otherIncome);
    const depositIncome = safeNumber(data.depositIncome);
    const rentExpense = safeNumber(data.rentExpense);
    const depositExpense = safeNumber(data.depositExpense);
    return rentIncome + waterIncome + electricIncome + otherIncome + depositIncome + rentExpense + depositExpense;
  };

  const effectiveTotal = incomeData ? getEffectiveTotal(incomeData) : 0;

  const incomeItems: IncomeItem[] = [
    {
      key: 'rentIncome',
      type: '租金收入',
      amount: safeNumber(incomeData?.rentIncome),
      formattedAmount: incomeData?.rentIncomeFormatted || '0.00',
      percentage: calculatePercentage(safeNumber(incomeData?.rentIncome), effectiveTotal),
    },
    {
      key: 'waterIncome',
      type: '水费收入',
      amount: safeNumber(incomeData?.waterIncome),
      formattedAmount: incomeData?.waterIncomeFormatted || '0.00',
      percentage: calculatePercentage(safeNumber(incomeData?.waterIncome), effectiveTotal),
    },
    {
      key: 'electricIncome',
      type: '电费收入',
      amount: safeNumber(incomeData?.electricIncome),
      formattedAmount: incomeData?.electricIncomeFormatted || '0.00',
      percentage: calculatePercentage(safeNumber(incomeData?.electricIncome), effectiveTotal),
    },
    {
      key: 'otherIncome',
      type: '其他收入',
      amount: safeNumber(incomeData?.otherIncome),
      formattedAmount: incomeData?.otherIncomeFormatted || '0.00',
      percentage: calculatePercentage(safeNumber(incomeData?.otherIncome), effectiveTotal),
    },
    {
      key: 'depositIncome',
      type: '押金收入',
      amount: safeNumber(incomeData?.depositIncome),
      formattedAmount: incomeData?.depositIncomeFormatted || '0.00',
      percentage: calculatePercentage(safeNumber(incomeData?.depositIncome), effectiveTotal),
    },
    {
      key: 'rentExpense',
      type: '租金支出',
      amount: safeNumber(incomeData?.rentExpense),
      formattedAmount: incomeData?.rentExpenseFormatted || '0.00',
      percentage: calculatePercentage(safeNumber(incomeData?.rentExpense), effectiveTotal),
      isExpense: true,
    },
    {
      key: 'depositExpense',
      type: '押金支出',
      amount: safeNumber(incomeData?.depositExpense),
      formattedAmount: incomeData?.depositExpenseFormatted || '0.00',
      percentage: calculatePercentage(safeNumber(incomeData?.depositExpense), effectiveTotal),
      isExpense: true,
    },
  ];

  const tableColumns = [
    {
      title: '收支类型',
      dataIndex: 'type',
      key: 'type',
      width: 120,
      render: (text: string, record: IncomeItem) => (
        <Tag color={record.isExpense ? 'red' : 'blue'}>
          {text}
        </Tag>
      ),
    },
    {
      title: '金额（元）',
      dataIndex: 'formattedAmount',
      key: 'formattedAmount',
      width: 140,
      render: (text: string, record: IncomeItem) => (
        <span style={{ color: record.isExpense ? '#cf1322' : 'inherit' }}>
          {record.isExpense ? '-' : ''}¥{text}
        </span>
      ),
    },
    {
      title: '占比',
      dataIndex: 'percentage',
      key: 'percentage',
      width: 100,
    },
  ];

  // 安全显示统计值
  const safeStatValue = (value: any): number => {
    const num = Number(value);
    return isNaN(num) ? 0 : num;
  };

  const netIncome = incomeData ? safeStatValue(incomeData.netIncome) : 0;
  const netIncomeColor = netIncome >= 0 ? '#3f8600' : '#cf1322';

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', flexWrap: 'wrap', gap: '16px', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1>收入查询</h1>
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: 16 }}>
          <Select
            placeholder="选择位置"
            onChange={handleLocationChange}
            allowClear={true}
            style={{ width: 200 }}
          >
            {(locations || []).map(location => (
              <Option key={location.id} value={location.id}>
                {location.shortName}
              </Option>
            ))}
          </Select>
          <MonthPicker
            placeholder="选择月份"
            onChange={handleMonthChange}
            allowClear={true}
            style={{ width: 200 }}
          />
          <Button
            type="primary"
            icon={<SearchOutlined />}
            onClick={() => fetchIncome(selectedMonth || undefined, selectedLocation || undefined)}
            loading={loading}
          >
            查询
          </Button>
          <Button
            icon={<DownloadOutlined />}
            onClick={handleExport}
          >
            导出
          </Button>
        </div>
      </div>

      {error && (
        <Alert
          message="错误"
          description={error}
          type="error"
          showIcon
          style={{ marginBottom: 16 }}
          action={
            <Button size="small" type="primary" onClick={() => fetchIncome(selectedMonth || undefined)}>
              重试
            </Button>
          }
        />
      )}

      <Spin spinning={loading} tip="加载中...">
        {incomeData && (
          <div>
            <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
              <Col xs={24} sm={12} lg={8}>
                <Card>
                  <Statistic
                    title="总收入"
                    value={safeStatValue(incomeData.totalIncome) / 100}
                    precision={2}
                    valueStyle={{ color: '#3f8600', fontSize: '28px' }}
                    prefix={<DollarOutlined />}
                    suffix="元"
                  />
                </Card>
              </Col>
              <Col xs={24} sm={12} lg={8}>
                <Card>
                  <Statistic
                    title="总支出"
                    value={safeStatValue(incomeData.totalExpense) / 100}
                    precision={2}
                    valueStyle={{ color: '#cf1322', fontSize: '28px' }}
                    prefix={<DollarOutlined />}
                    suffix="元"
                  />
                </Card>
              </Col>
              <Col xs={24} sm={12} lg={8}>
                <Card>
                  <Statistic
                    title="净收入"
                    value={netIncome / 100}
                    precision={2}
                    valueStyle={{ color: netIncomeColor, fontSize: '28px' }}
                    prefix={<DollarOutlined />}
                    suffix="元"
                  />
                </Card>
              </Col>
            </Row>

            <Card title="收支明细">
              <Table
                dataSource={incomeItems}
                columns={tableColumns}
                pagination={false}
                bordered={false}
              />
            </Card>
          </div>
        )}
      </Spin>
    </div>
  );
};

export default Income;
