import React, { createContext, useContext, useReducer, useEffect, useCallback } from 'react';
import { useAuth } from './AuthContext';
import { telemetryAPI } from '../services/api';

const DeviceContext = createContext();

const initialState = {
  devices: [],
  selectedDevice: null,
  loading: false,
  error: null,
};

const deviceReducer = (state, action) => {
  switch (action.type) {
    case 'SET_LOADING':
      return { ...state, loading: action.payload };
    case 'SET_DEVICES':
      return { ...state, devices: action.payload, loading: false };
    case 'SET_SELECTED_DEVICE':
      return { ...state, selectedDevice: action.payload };
    case 'SET_ERROR':
      return { ...state, error: action.payload, loading: false };
    case 'CLEAR_ERROR':
      return { ...state, error: null };
    default:
      return state;
  }
};

export const DeviceProvider = ({ children }) => {
  const [state, dispatch] = useReducer(deviceReducer, initialState);
  const { user, isAuthenticated } = useAuth();

  const fetchDevices = useCallback(async () => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      const response = await telemetryAPI.getDevices();
      dispatch({ type: 'SET_DEVICES', payload: response.data });

      // Auto-select first device if none selected
      if (!state.selectedDevice && response.data.length > 0) {
        dispatch({ type: 'SET_SELECTED_DEVICE', payload: response.data[0] });
      }
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: error.response?.data?.error || 'Failed to fetch devices' });
    }
  }, [state.selectedDevice]);

  useEffect(() => {
    if (isAuthenticated && user) {
      fetchDevices();
    }
  }, [isAuthenticated, user, fetchDevices]);

  const selectDevice = (device) => {
    dispatch({ type: 'SET_SELECTED_DEVICE', payload: device });
  };

  const clearError = () => {
    dispatch({ type: 'CLEAR_ERROR' });
  };

  const value = {
    ...state,
    fetchDevices,
    selectDevice,
    clearError,
  };

  return (
    <DeviceContext.Provider value={value}>
      {children}
    </DeviceContext.Provider>
  );
};

export const useDevice = () => {
  const context = useContext(DeviceContext);
  if (!context) {
    throw new Error('useDevice must be used within a DeviceProvider');
  }
  return context;
};
