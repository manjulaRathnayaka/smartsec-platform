import { useState, useEffect } from 'react';

export const useDeviceId = () => {
  const [deviceId, setDeviceId] = useState(null);

  useEffect(() => {
    // Try to get device ID from localStorage first
    let storedDeviceId = localStorage.getItem('deviceId');

    if (!storedDeviceId) {
      // Generate a simple device ID based on browser fingerprint
      const canvas = document.createElement('canvas');
      const ctx = canvas.getContext('2d');
      ctx.textBaseline = 'top';
      ctx.font = '14px Arial';
      ctx.fillText('Device fingerprint', 2, 2);

      const fingerprint = [
        navigator.userAgent,
        navigator.language,
        screen.width,
        screen.height,
        new Date().getTimezoneOffset(),
        canvas.toDataURL()
      ].join('|');

      // Simple hash function
      let hash = 0;
      for (let i = 0; i < fingerprint.length; i++) {
        const char = fingerprint.charCodeAt(i);
        hash = ((hash << 5) - hash) + char;
        hash = hash & hash; // Convert to 32-bit integer
      }

      storedDeviceId = `device_${Math.abs(hash).toString(16)}`;
      localStorage.setItem('deviceId', storedDeviceId);
    }

    setDeviceId(storedDeviceId);
  }, []);

  return deviceId;
};
