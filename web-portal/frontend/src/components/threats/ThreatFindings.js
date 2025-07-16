import React, { useState } from 'react';
import { useQuery } from 'react-query';
import { useDevice } from '../../contexts/DeviceContext';
import { telemetryAPI } from '../../services/api';
import {
  ExclamationTriangleIcon,
  ShieldExclamationIcon,
  FunnelIcon,
  CalendarIcon
} from '@heroicons/react/24/outline';

const ThreatFindings = () => {
  const { selectedDevice } = useDevice();
  const [severityFilter, setSeverityFilter] = useState('all');
  const [timeRange, setTimeRange] = useState('24h');

  const { data: threatData, isLoading } = useQuery(
    ['device-threats', selectedDevice?.id, timeRange],
    () => telemetryAPI.getThreats(selectedDevice?.id),
    {
      enabled: !!selectedDevice,
      refetchInterval: 30000,
    }
  );

  const threats = threatData?.data || [];
  const filteredThreats = threats.filter(threat =>
    severityFilter === 'all' || threat.severity === severityFilter
  );

  const getSeverityColor = (severity) => {
    switch (severity) {
      case 'critical': return 'bg-red-100 text-red-800 border-red-200';
      case 'high': return 'bg-orange-100 text-orange-800 border-orange-200';
      case 'medium': return 'bg-yellow-100 text-yellow-800 border-yellow-200';
      case 'low': return 'bg-green-100 text-green-800 border-green-200';
      default: return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  const getSeverityIcon = (severity) => {
    switch (severity) {
      case 'critical':
      case 'high':
        return <ExclamationTriangleIcon className="h-5 w-5 text-red-500" />;
      case 'medium':
        return <ShieldExclamationIcon className="h-5 w-5 text-yellow-500" />;
      default:
        return <ShieldExclamationIcon className="h-5 w-5 text-green-500" />;
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="md:flex md:items-center md:justify-between">
        <div className="flex-1 min-w-0">
          <h2 className="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
            Threat Findings
          </h2>
          <p className="mt-1 text-sm text-gray-500">
            Security threats and anomalies detected on your devices
          </p>
        </div>
      </div>

      {/* Filters */}
      <div className="bg-white shadow rounded-lg p-6">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-4 sm:space-y-0">
          <div className="flex items-center space-x-4">
            <div className="relative">
              <FunnelIcon className="h-5 w-5 text-gray-400 absolute left-3 top-3" />
              <select
                value={severityFilter}
                onChange={(e) => setSeverityFilter(e.target.value)}
                className="pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="all">All Severities</option>
                <option value="critical">Critical</option>
                <option value="high">High</option>
                <option value="medium">Medium</option>
                <option value="low">Low</option>
              </select>
            </div>
          </div>

          <div className="flex items-center space-x-4">
            <div className="relative">
              <CalendarIcon className="h-5 w-5 text-gray-400 absolute left-3 top-3" />
              <select
                value={timeRange}
                onChange={(e) => setTimeRange(e.target.value)}
                className="pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="1h">Last hour</option>
                <option value="24h">Last 24 hours</option>
                <option value="7d">Last 7 days</option>
                <option value="30d">Last 30 days</option>
              </select>
            </div>
          </div>
        </div>
      </div>

      {/* Threat Summary */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        {['critical', 'high', 'medium', 'low'].map((severity) => {
          const count = threats.filter(t => t.severity === severity).length;
          return (
            <div key={severity} className="bg-white overflow-hidden shadow rounded-lg">
              <div className="p-5">
                <div className="flex items-center">
                  <div className="flex-shrink-0">
                    {getSeverityIcon(severity)}
                  </div>
                  <div className="ml-5 w-0 flex-1">
                    <dl>
                      <dt className="text-sm font-medium text-gray-500 truncate capitalize">
                        {severity} Threats
                      </dt>
                      <dd className="text-lg font-medium text-gray-900">{count}</dd>
                    </dl>
                  </div>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Threats List */}
      <div className="bg-white shadow overflow-hidden sm:rounded-md">
        {isLoading ? (
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          </div>
        ) : filteredThreats.length === 0 ? (
          <div className="text-center py-12">
            <ShieldExclamationIcon className="mx-auto h-12 w-12 text-gray-400" />
            <h3 className="mt-2 text-sm font-medium text-gray-900">No threats found</h3>
            <p className="mt-1 text-sm text-gray-500">
              Your device appears to be secure. Keep monitoring for any changes.
            </p>
          </div>
        ) : (
          <ul className="divide-y divide-gray-200">
            {filteredThreats.map((threat) => (
              <li key={threat.id} className="px-6 py-4 hover:bg-gray-50">
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    <div className="flex-shrink-0">
                      {getSeverityIcon(threat.severity)}
                    </div>
                    <div className="ml-4">
                      <div className="flex items-center">
                        <p className="text-sm font-medium text-gray-900">{threat.title}</p>
                        <span className={`ml-2 inline-flex px-2 py-1 text-xs font-semibold rounded-full border ${getSeverityColor(threat.severity)}`}>
                          {threat.severity}
                        </span>
                      </div>
                      <p className="text-sm text-gray-500">{threat.description}</p>
                      <div className="mt-2 flex items-center text-sm text-gray-500">
                        <span>Device: {threat.device_hostname}</span>
                        <span className="mx-2">•</span>
                        <span>Process: {threat.process_name}</span>
                        <span className="mx-2">•</span>
                        <span>{new Date(threat.detected_at).toLocaleString()}</span>
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center space-x-2">
                    <button className="text-blue-600 hover:text-blue-900 text-sm font-medium">
                      Investigate
                    </button>
                    <button className="text-gray-600 hover:text-gray-900 text-sm font-medium">
                      Dismiss
                    </button>
                  </div>
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
};

export default ThreatFindings;
