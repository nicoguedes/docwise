import { useState, useEffect, useCallback, useRef } from 'react';
import { Document } from '../types';
import * as api from '../api/client';

export function useDocuments() {
  const [documents, setDocuments] = useState<Document[]>([]);
  const [loading, setLoading] = useState(true);
  const intervalRef = useRef<ReturnType<typeof setInterval>>();

  const refresh = useCallback(async () => {
    try {
      const docs = await api.listDocuments();
      setDocuments(docs);
    } catch (err) {
      console.error('Failed to fetch documents:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();

    // Poll for status updates while any document is processing
    intervalRef.current = setInterval(() => {
      refresh();
    }, 3000);

    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current);
    };
  }, [refresh]);

  const upload = useCallback(async (file: File) => {
    const doc = await api.uploadDocument(file);
    setDocuments(prev => [doc, ...prev]);
    return doc;
  }, []);

  const remove = useCallback(async (id: string) => {
    await api.deleteDocument(id);
    setDocuments(prev => prev.filter(d => d.id !== id));
  }, []);

  return { documents, loading, upload, remove, refresh };
}
