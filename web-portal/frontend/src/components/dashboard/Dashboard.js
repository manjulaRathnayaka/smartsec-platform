import React, { useState, useEffect } from 'react';
import { useQuery } from 'react-query';
import { useAuth } from '../../contexts/AuthContext';
import { useDevice } from '../../contexts/DeviceContext';
import { telemetryAPI, mcpAPI } from '../../services/api';
import {
  ComputerDesktopIcon,
  ExclamationTriangleIcon,
  HeartIcon,
  EyeIcon,
  ArrowTrendingUpIcon,
  ArrowTrendingDownIcon
} from '@heroicons/react/24/outline';

const Dashboard = () => {
  const { user } = useAuth();
  const { selectedDevice } = useDevice();
  const [timeRange, setTimeRange] = useState('24h');

  const { data: dashboardData, isLoading } = useQuery(
    ['dashboard', selectedDevice?.id, timeRange],
    () => telemetryAPI.getDashboardData(selectedDevice?.id, timeRange),
    {
      enabled: !!selectedDevice,
      refetchInterval: 30000, // Refresh every 30 seconds
    }
  );

  const stats = dashboardData?.data || {
    activeProcesses: 0,
    runningContainers: 0,
    threatsDetected: 0,
    healthScore: 0,
    cpuUsage: 0,
    memoryUsage: 0,
    diskUsage: 0,
    networkActivity: 0,
  };

  const StatCard = ({ title, value, icon: Icon, trend, trendValue, color = 'blue' }) => (
    <div className="bg-white overflow-hidden shadow rounded-lg">
      <div className="p-5">
        <div className="flex items-center">
          <div className="flex-shrink-0">
            <Icon className={`h-6 w-6 text-${color}-600`} />
          </div>
          <div className="ml-5 w-0 flex-1">
            <dl>
              <dt className="text-sm font-medium text-gray-500 truncate">{title}</dt>
              <dd className="flex items-baseline">
                <div className="text-2xl font-semibold text-gray-900">{value}</div>
                {trend && (
                  <div className={`ml-2 flex items-baseline text-sm font-semibold ${
                    trend === 'up' ? 'text-green-600' : 'text-red-600'
                  }`}>
                    {trend === 'up' ? (
                      <ArrowTrendingUpIcon className="h-4 w-4 flex-shrink-0" />
                    ) : (
                      <ArrowTrendingDownIcon className="h-4 w-4 flex-shrink-0" />
                    )}
                    <span className="sr-only">{trend === 'up' ? 'Increased' : 'Decreased'} by</span>
                    {trendValue}
                  </div>
                )}
              </dd>
            </dl>
          </div>
        </div>
      </div>
    </div>
  );

  const HealthBar = ({ label, value, color = 'blue' }) => (
    <div className="mb-4">
      <div className="flex justify-between text-sm font-medium text-gray-900 mb-1">
        <span>{label}</span>
        <span>{value}%</span>
      </div>
      <div className="w-full bg-gray-200 rounded-full h-2">
        <div
          className={`bg-${color}-600 h-2 rounded-full transition-all duration-300`}
          style={{ width: `${value}%` }}
        />
      </div>
    </div>
  );

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="md:flex md:items-center md:justify-between">
        <div className="flex-1 min-w-0">
          <h2 className="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
            Security Dashboard
          </h2>
          <p className="mt-1 text-sm text-gray-500">
            Welcome back, {user?.name}. Here's your security overview.
          </p>
        </div>
        <div className="mt-4 flex md:mt-0 md:ml-4">
          <select
            value={timeRange}
            onChange={(e) => setTimeRange(e.target.value)}
            className="rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
          >
            <option value="1h">Last hour</option>
            <option value="24h">Last 24 hours</option>
            <option value="7d">Last 7 days</option>
            <option value="30d">Last 30 days</option>
          </select>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard
          title="Active Processes"
          value={stats.activeProcesses}
          icon={ComputerDesktopIcon}
          trend="up"
          trendValue="+12%"
        />
        <StatCard
          title="Running Containers"
          value={stats.runningContainers}
          icon={EyeIcon}
          trend="down"
          trendValue="-3%"
          color="green"
        />
        <StatCard
          title="Threats Detected"
          value={stats.threatsDetected}
          icon={ExclamationTriangleIcon}
          trend={stats.threatsDetected > 0 ? 'up' : 'down'}
          trendValue={stats.threatsDetected > 0 ? '+2' : '0'}
          color="red"
        />
        <StatCard
          title="Health Score"
          value={`${stats.healthScore}%`}
          icon={HeartIcon}
          trend="up"
          trendValue="+5%"
          color="green"
        />
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* System Health */}
        <div className="bg-white shadow rounded-lg p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">System Health</h3>
          <HealthBar label="CPU Usage" value={stats.cpuUsage} color="blue" />
          <HealthBar label="Memory Usage" value={stats.memoryUsage} color="yellow" />
          <HealthBar label="Disk Usage" value={stats.diskUsage} color="red" />
          <HealthBar label="Network Activity" value={stats.networkActivity} color="green" />
        </div>

        {/* Recent Activity */}
        <div className="bg-white shadow rounded-lg p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Recent Activity</h3>
          <div className="flow-root">
            <ul className="-mb-8">
              {[
                { id: 1, type: 'process', message: 'New process started: nginx', time: '2 minutes ago' },
                { id: 2, type: 'container', message: 'Container deployed: web-app:latest', time: '15 minutes ago' },
                { id: 3, type: 'threat', message: 'Suspicious activity detected', time: '1 hour ago' },
                { id: 4, type: 'system', message: 'System health check completed', time: '2 hours ago' },
              ].map((item, itemIdx) => (
                <li key={item.id}>
                  <div className="relative pb-8">
                    {itemIdx !== 3 && (
                      <span className="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200" />
                    )}
                    <div className="relative flex space-x-3">
                      <div>
                        <span className={`h-8 w-8 rounded-full flex items-center justify-center ring-8 ring-white ${
                          item.type === 'threat' ? 'bg-red-500' :
                          item.type === 'container' ? 'bg-blue-500' : 'bg-green-500'
                        }`}>
                          {item.type === 'process' && <ComputerDesktopIcon className="h-5 w-5 text-white" />}
                          {item.type === 'container' && <EyeIcon className="h-5 w-5 text-white" />}
                          {item.type === 'threat' && <ExclamationTriangleIcon className="h-5 w-5 text-white" />}
                          {item.type === 'system' && <HeartIcon className="h-5 w-5 text-white" />}
                        </span>
                      </div>
                      <div className="min-w-0 flex-1 pt-1.5 flex justify-between space-x-4">
                        <div>
                          <p className="text-sm text-gray-500">{item.message}</p>
                        </div>
                        <div className="text-right text-sm whitespace-nowrap text-gray-500">
                          {item.time}
                        </div>
                      </div>
                    </div>
                  </div>
                </li>
              ))}
            </ul>
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="bg-white shadow rounded-lg p-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4">Quick Actions</h3>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <button className="flex items-center justify-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50">
            <ComputerDesktopIcon className="h-5 w-5 mr-2" />
            View Devices
          </button>
          <button className="flex items-center justify-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50">
            <ExclamationTriangleIcon className="h-5 w-5 mr-2" />
            Check Threats
          </button>
          <button className="flex items-center justify-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50">
            <HeartIcon className="h-5 w-5 mr-2" />
            System Health
          </button>
          <button className="flex items-center justify-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50">
            <EyeIcon className="h-5 w-5 mr-2" />
            AI Assistant
          </button>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
