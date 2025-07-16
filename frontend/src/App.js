import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import { DeviceProvider } from './contexts/DeviceContext';
import ErrorBoundary from './components/common/ErrorBoundary';
import ProtectedRoute from './components/common/ProtectedRoute';
import Layout from './components/layout/Layout';
import Login from './components/auth/Login';
import AuthCallback from './components/auth/AuthCallback';
import Dashboard from './components/dashboard/Dashboard';
import DeviceActivity from './components/devices/DeviceActivity';
import ThreatFindings from './components/threats/ThreatFindings';
import SystemHealth from './components/health/SystemHealth';
import AdminDashboard from './components/admin/AdminDashboard';
import FleetActivity from './components/admin/FleetActivity';
import AIAssistant from './components/ai/AIAssistant';
import Profile from './components/profile/Profile';

function App() {
  return (
    <ErrorBoundary>
      <AuthProvider>
        <DeviceProvider>
          <div className="App">
            <Routes>
              <Route path="/login" element={<Login />} />
              <Route path="/auth/callback" element={<AuthCallback />} />
              <Route path="/" element={<ProtectedRoute><Layout /></ProtectedRoute>}>
                <Route index element={<Navigate to="/dashboard" replace />} />
                <Route path="dashboard" element={<Dashboard />} />
                <Route path="devices" element={<DeviceActivity />} />
                <Route path="threats" element={<ThreatFindings />} />
                <Route path="health" element={<SystemHealth />} />
                <Route path="ai" element={<AIAssistant />} />
                <Route path="profile" element={<Profile />} />

                {/* Admin routes */}
                <Route path="admin" element={<ProtectedRoute requireAdmin><AdminDashboard /></ProtectedRoute>} />
                <Route path="admin/fleet" element={<ProtectedRoute requireAdmin><FleetActivity /></ProtectedRoute>} />
              </Route>
            </Routes>
          </div>
        </DeviceProvider>
      </AuthProvider>
    </ErrorBoundary>
  );
}

export default App;
