import { useEffect } from 'react';
import { wsService } from '../services/websocket.service';

export const useWebSocket = (event: string, callback: (data: any) => void) => {
  useEffect(() => {
    wsService.on(event, callback);

    return () => {
      wsService.off(event, callback);
    };
  }, [event, callback]);
};
