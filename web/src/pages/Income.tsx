import React, { useState, useEffect } from 'react';
import { Card, DatePicker, Statistic, Table, Button, Spin, Alert, message, Row, Col, Tag } from 'antd';
import { DollarOutlined, DownloadOutlined, SearchOutlined } from '@ant-design/icons';
import { incomeApi } from '../api/income';

const { MonthPicker } = DatePicker;

interface IncomeItem {
  key: string;
  type: string;
  amount: number;
  formattedAmount: string;
  percentage: string;
}

// 安全计算占比函数
const calculatePercentage = (value: number, total: number): string => {
  if (!total || total === 0 || isNaN(total)) {
    return '0%';
  }
  // 使用绝对值计算占比（押金支出是负数）
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

  const fetchIncome = async (month?: string) => {
    setLoading(true);
    setError(null);
    try {
      const data = await incomeApi.getReport(month);
      setIncomeData(data);
      message.success('获取收入数据成功');
    } catch (err) {
      setError('获取收入数据失败');
      console.error('Failed to fetch income:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchIncome();
  }, []);

  const handleMonthChange = (date: any) => {
    if (date) {
      const monthStr = date.format('YYYY-MM');
      setSelectedMonth(monthStr);
      fetchIncome(monthStr);
    } else {
      setSelectedMonth(null);
      fetchIncome();
    }
  };

  const handleExport = () => {
    message.info('导出功能开发中...');
  };

  // 安全获取数值，避免 NaN
  const safeNumber = (value: any, defaultValue: number = 0): number => {
    const num = Number(value);
    return isNaN(num) ? defaultValue : num;
  };

  // 计算用于占比计算的有效总金额（押金支出取绝对值）
  const getEffectiveTotal = (data: any): number => {
    if (!data) return 0;
    const rent = safeNumber(data.rentTotal);
    const water = safeNumber(data.waterTotal);
    const electric = safeNumber(data.electricTotal);
    const other = safeNumber(data.otherTotal);
    const depositIncome = safeNumber(data.depositIncome);
    const depositExpense = safeNumber(data.depositExpense);
    // 押金支出取绝对值，其他保持原样
    return rent + water + electric + other + depositIncome + Math.abs(depositExpense);
  };

  const effectiveTotal = incomeData ? getEffectiveTotal(incomeData) : 0;

  const incomeItems: IncomeItem[] = [
    {
      key: 'rent',
      type: '租金',
      amount: safeNumber(incomeData?.rentTotal),
      formattedAmount: incomeData?.rentFormatted || '0.00',
      percentage: calculatePercentage(safeNumber(incomeData?.rentTotal), effectiveTotal),
    },
    {
      key: 'water',
      type: '水费',
      amount: safeNumber(incomeData?.waterTotal),
      formattedAmount: incomeData?.waterFormatted || '0.00',
      percentage: calculatePercentage(safeNumber(incomeData?.waterTotal), effectiveTotal),
    },
    {
      key: 'electric',
      type: '电费',
      amount: safeNumber(incomeData?.electricTotal),
      formattedAmount: incomeData?.electricFormatted || '0.00',
      percentage: calculatePercentage(safeNumber(incomeData?.electricTotal), effectiveTotal),
    },
    {
      key: 'other',
      type: '其他',
      amount: safeNumber(incomeData?.otherTotal),
      formattedAmount: incomeData?.otherFormatted || '0.00',
      percentage: calculatePercentage(safeNumber(incomeData?.otherTotal), effectiveTotal),
    },
    {
      key: 'depositIncome',
      type: '押金收入',
      amount: safeNumber(incomeData?.depositIncome),
      formattedAmount: incomeData?.depositIncomeFormatted || '0.00',
      percentage: calculatePercentage(safeNumber(incomeData?.depositIncome), effectiveTotal),
    },
    {
      key: 'depositExpense',
      type: '押金支出',
      amount: safeNumber(incomeData?.depositExpense),
      formattedAmount: incomeData?.depositExpenseFormatted || '0.00',
      percentage: calculatePercentage(safeNumber(incomeData?.depositExpense), effectiveTotal),
    },
  ];

  const tableColumns = [
    {
      title: '收入类型',
      dataIndex: 'type',
      key: 'type',
      width: 120,
      render: (text: string, record: IncomeItem) => (
        <Tag color={record.key === 'depositExpense' ? 'red' : 'blue'}>
          {text}
        </Tag>
      ),
    },
    {
      title: '金额（元）',
      dataIndex: 'formattedAmount',
      key: 'formattedAmount',
      width: 120,
      render: (text: string, record: IncomeItem) => (
        <span style={{ color: record.key === 'depositExpense' ? '#cf1322' : 'inherit' }}>
          {record.key === 'depositExpense' ? '-' : ''}¥{text}
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

  const monthTitle = selectedMonth
    ? `${selectedMonth} 收入汇总`
    : '本月收入汇总';

  // 安全显示统计值
  const safeStatValue = (value: any): number => {
    const num = Number(value);
    return isNaN(num) ? 0 : num;
  };

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', flexWrap: 'wrap', gap: '16px', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1>收入查询</h1>
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: 16 }}>
          <MonthPicker
            placeholder="选择月份"
            onChange={handleMonthChange}
            allowClear={true}
            style={{ width: 200 }}
          />
          <Button
            type="primary"
            icon={<SearchOutlined />}
            onClick={() => fetchIncome(selectedMonth || undefined)}
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
                    title={monthTitle}
                    value={safeStatValue(incomeData.total) / 100}
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
                    title="押金收入"
                    value={safeStatValue(incomeData.depositIncome) / 100}
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
                    title="押金支出"
                    value={safeStatValue(incomeData.depositExpense) / 100}
                    precision={2}
                    valueStyle={{ color: '#cf1322', fontSize: '28px' }}
                    prefix={<DollarOutlined />}
                    suffix="元"
                  />
                </Card>
              </Col>
            </Row>

            <Card title="收入明细">
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
