import React, { useState } from 'react';
import { useQuery } from 'react-query';
import { useDevice } from '../../contexts/DeviceContext';
import { telemetryAPI } from '../../services/api';
import {
  ComputerDesktopIcon,
  EyeIcon,
  FunnelIcon,
  ArrowsUpDownIcon,
  PlayIcon,
  StopIcon,
  ClockIcon
} from '@heroicons/react/24/outline';

const DeviceActivity = () => {
  const { selectedDevice } = useDevice();
  const [activeTab, setActiveTab] = useState('processes');
  const [filter, setFilter] = useState('');
  const [sortBy, setSortBy] = useState('start_time');
  const [sortOrder, setSortOrder] = useState('desc');

  const { data: processData, isLoading: processLoading } = useQuery(
    ['device-processes', selectedDevice?.id],
    () => telemetryAPI.getDeviceActivity(selectedDevice?.id),
    {
      enabled: !!selectedDevice && activeTab === 'processes',
      refetchInterval: 10000,
    }
  );

  const { data: containerData, isLoading: containerLoading } = useQuery(
    ['device-containers', selectedDevice?.id],
    () => telemetryAPI.getContainers(selectedDevice?.id),
    {
      enabled: !!selectedDevice && activeTab === 'containers',
      refetchInterval: 10000,
    }
  );

  const processes = processData?.data || [];
  const containers = containerData?.data || [];

  const filteredProcesses = processes.filter(process =>
    process.name.toLowerCase().includes(filter.toLowerCase()) ||
    process.cmdline.toLowerCase().includes(filter.toLowerCase()) ||
    process.username.toLowerCase().includes(filter.toLowerCase())
  );

  const filteredContainers = containers.filter(container =>
    container.image_name.toLowerCase().includes(filter.toLowerCase()) ||
    container.container_id.toLowerCase().includes(filter.toLowerCase()) ||
    container.user.toLowerCase().includes(filter.toLowerCase())
  );

  const ProcessTable = () => (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Process
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              PID
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              User
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Command
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Started
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Status
            </th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {filteredProcesses.map((process) => (
            <tr key={process.pid} className="hover:bg-gray-50">
              <td className="px-6 py-4 whitespace-nowrap">
                <div className="flex items-center">
                  <PlayIcon className="h-5 w-5 text-green-500 mr-2" />
                  <div>
                    <div className="text-sm font-medium text-gray-900">{process.name}</div>
                    <div className="text-sm text-gray-500">{process.exe_path}</div>
                  </div>
                </div>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {process.pid}
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                  process.username === 'root'
                    ? 'bg-red-100 text-red-800'
                    : 'bg-blue-100 text-blue-800'
                }`}>
                  {process.username}
                </span>
              </td>
              <td className="px-6 py-4 text-sm text-gray-900 max-w-md truncate">
                {process.cmdline}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {new Date(process.start_time).toLocaleString()}
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <span className="inline-flex px-2 py-1 text-xs font-semibold rounded-full bg-green-100 text-green-800">
                  Running
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );

  const ContainerTable = () => (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Container
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Image
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              User
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Ports
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Created
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Status
            </th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {filteredContainers.map((container) => (
            <tr key={container.container_id} className="hover:bg-gray-50">
              <td className="px-6 py-4 whitespace-nowrap">
                <div className="flex items-center">
                  <EyeIcon className="h-5 w-5 text-blue-500 mr-2" />
                  <div>
                    <div className="text-sm font-medium text-gray-900">
                      {container.container_id.substring(0, 12)}
                    </div>
                    <div className="text-sm text-gray-500">{container.name}</div>
                  </div>
                </div>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {container.image_name}
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                  container.user === 'root'
                    ? 'bg-red-100 text-red-800'
                    : 'bg-blue-100 text-blue-800'
                }`}>
                  {container.user}
                </span>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {container.ports && container.ports.length > 0
                  ? container.ports.join(', ')
                  : 'None'
                }
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {new Date(container.created_at).toLocaleString()}
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                  container.status === 'running'
                    ? 'bg-green-100 text-green-800'
                    : 'bg-yellow-100 text-yellow-800'
                }`}>
                  {container.status}
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );

  const isLoading = processLoading || containerLoading;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="md:flex md:items-center md:justify-between">
        <div className="flex-1 min-w-0">
          <h2 className="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
            Device Activity
          </h2>
          <p className="mt-1 text-sm text-gray-500">
            {selectedDevice
              ? `Monitoring ${selectedDevice.hostname} (${selectedDevice.mac_address})`
              : 'No device selected'
            }
          </p>
        </div>
      </div>

      {/* Tabs */}
      <div className="border-b border-gray-200">
        <nav className="-mb-px flex space-x-8">
          <button
            onClick={() => setActiveTab('processes')}
            className={`whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'processes'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            <ComputerDesktopIcon className="h-5 w-5 inline-block mr-2" />
            Processes
          </button>
          <button
            onClick={() => setActiveTab('containers')}
            className={`whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'containers'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            <EyeIcon className="h-5 w-5 inline-block mr-2" />
            Containers
          </button>
        </nav>
      </div>

      {/* Filters */}
      <div className="bg-white shadow rounded-lg p-6">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-4 sm:space-y-0">
          <div className="flex items-center space-x-4">
            <div className="relative">
              <FunnelIcon className="h-5 w-5 text-gray-400 absolute left-3 top-3" />
              <input
                type="text"
                placeholder={`Filter ${activeTab}...`}
                value={filter}
                onChange={(e) => setFilter(e.target.value)}
                className="pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>

          <div className="flex items-center space-x-4">
            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value)}
              className="rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
            >
              <option value="start_time">Start Time</option>
              <option value="name">Name</option>
              <option value="username">User</option>
              <option value="pid">PID</option>
            </select>

            <button
              onClick={() => setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc')}
              className="p-2 border border-gray-300 rounded-md hover:bg-gray-50"
            >
              <ArrowsUpDownIcon className="h-5 w-5" />
            </button>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="bg-white shadow overflow-hidden sm:rounded-md">
        {isLoading ? (
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          </div>
        ) : (
          <>
            {activeTab === 'processes' && (
              <div>
                <div className="px-6 py-3 border-b border-gray-200 bg-gray-50">
                  <p className="text-sm text-gray-700">
                    {filteredProcesses.length} processes found
                  </p>
                </div>
                <ProcessTable />
              </div>
            )}

            {activeTab === 'containers' && (
              <div>
                <div className="px-6 py-3 border-b border-gray-200 bg-gray-50">
                  <p className="text-sm text-gray-700">
                    {filteredContainers.length} containers found
                  </p>
                </div>
                <ContainerTable />
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

export default DeviceActivity;
