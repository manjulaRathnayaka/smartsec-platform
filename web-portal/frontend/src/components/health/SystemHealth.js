import React from 'react';
import { useQuery } from 'react-query';
import { useDevice } from '../../contexts/DeviceContext';
import { telemetryAPI } from '../../services/api';
import { HeartIcon, CpuChipIcon, CircleStackIcon, GlobeAltIcon } from '@heroicons/react/24/outline';

const SystemHealth = () => {
  const { selectedDevice } = useDevice();

  const { data: healthData, isLoading } = useQuery(
    ['device-health', selectedDevice?.id],
    () => telemetryAPI.getSystemHealth(selectedDevice?.id),
    {
      enabled: !!selectedDevice,
      refetchInterval: 10000,
    }
  );

  const health = healthData?.data || {
    cpu_usage: 0,
    memory_usage: 0,
    disk_usage: 0,
    network_usage: 0,
    uptime: 0,
    overall_score: 0,
  };

  const HealthMetric = ({ icon: Icon, label, value, unit, color }) => (
    <div className="bg-white overflow-hidden shadow rounded-lg">
      <div className="p-5">
        <div className="flex items-center">
          <div className="flex-shrink-0">
            <Icon className={`h-6 w-6 text-${color}-600`} />
          </div>
          <div className="ml-5 w-0 flex-1">
            <dl>
              <dt className="text-sm font-medium text-gray-500 truncate">{label}</dt>
              <dd className="text-2xl font-semibold text-gray-900">
                {value}{unit}
              </dd>
            </dl>
          </div>
        </div>
        <div className="mt-3">
          <div className="bg-gray-200 rounded-full h-2">
            <div
              className={`bg-${color}-600 h-2 rounded-full transition-all duration-300`}
              style={{ width: `${value}%` }}
            />
          </div>
        </div>
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
            System Health
          </h2>
          <p className="mt-1 text-sm text-gray-500">
            Real-time performance metrics and system status
          </p>
        </div>
      </div>

      {/* Overall Health Score */}
      <div className="bg-white shadow rounded-lg p-6">
        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-lg font-medium text-gray-900">Overall Health Score</h3>
            <p className="text-sm text-gray-500">System performance indicator</p>
          </div>
          <div className="text-right">
            <div className="text-3xl font-bold text-green-600">{health.overall_score}%</div>
            <div className="text-sm text-gray-500">Excellent</div>
          </div>
        </div>
      </div>

      {/* Metrics Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <HealthMetric
          icon={CpuChipIcon}
          label="CPU Usage"
          value={health.cpu_usage}
          unit="%"
          color="blue"
        />
        <HealthMetric
          icon={CircleStackIcon}
          label="Memory Usage"
          value={health.memory_usage}
          unit="%"
          color="yellow"
        />
        <HealthMetric
          icon={CircleStackIcon}
          label="Disk Usage"
          value={health.disk_usage}
          unit="%"
          color="red"
        />
        <HealthMetric
          icon={GlobeAltIcon}
          label="Network Usage"
          value={health.network_usage}
          unit="%"
          color="green"
        />
      </div>

      {/* Additional Info */}
      <div className="bg-white shadow rounded-lg p-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4">System Information</h3>
        <dl className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <dt className="text-sm font-medium text-gray-500">Uptime</dt>
            <dd className="text-sm text-gray-900">{Math.floor(health.uptime / 3600)} hours</dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500">Last Updated</dt>
            <dd className="text-sm text-gray-900">{new Date().toLocaleString()}</dd>
          </div>
        </dl>
      </div>
    </div>
  );
};

export default SystemHealth;
