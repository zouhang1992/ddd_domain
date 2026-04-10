import React from 'react';
import { Spin, Typography } from 'antd';

const { Title, Text } = Typography;

interface LoadingPageProps {
  title?: string;
  description?: string;
  icon?: React.ReactNode;
}

export const LoadingPage: React.FC<LoadingPageProps> = ({
  title = '加载中...',
  description = '正在为您准备页面，请稍候',
  icon,
}) => {
  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'center',
      alignItems: 'center',
      height: '100vh',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    }}>
      <div style={{
        background: 'white',
        padding: '48px 64px',
        borderRadius: '16px',
        boxShadow: '0 20px 60px rgba(0, 0, 0, 0.3)',
        textAlign: 'center',
        maxWidth: '400px',
      }}>
        {icon || (
          <Spin
            size="large"
            style={{
              fontSize: '48px',
              marginBottom: '24px',
            }}
          />
        )}
        <Title
          level={3}
          style={{
            marginBottom: '8px',
            color: '#1a1a1a',
          }}
        >
          {title}
        </Title>
        <Text
          type="secondary"
          style={{
            fontSize: '14px',
          }}
        >
          {description}
        </Text>
      </div>
    </div>
  );
};

export const RedirectingPage: React.FC<{ to?: string }> = ({ to }) => {
  return (
    <LoadingPage
      title="正在跳转..."
      description={to ? `即将跳转到 ${to}` : '正在为您重定向到目标页面'}
    />
  );
};

export const AuthLoadingPage: React.FC = () => {
  return (
    <LoadingPage
      title="验证登录状态..."
      description="正在检查您的登录信息，请稍候"
    />
  );
};

export default LoadingPage;
