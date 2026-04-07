import { useEffect, useState } from 'react';
import { initializeApp } from 'firebase/app';
import { getFirestore, doc, onSnapshot } from 'firebase/firestore';

// Environment configurations bound directly ensuring global deployment routing
const firebaseConfig = {
  projectId: process.env.REACT_APP_FIREBASE_PROJECT_ID,
  // ... standardized firebase configs
};

const app = initializeApp(firebaseConfig);
const db = getFirestore(app);

export interface Heatmap {
  zone_id: string;
  density_level: number;
  timestamp: string;
}

export interface WaitTime {
  amenity_id: string;
  wait_time_minutes: number;
}

/**
 * CLIENT SDK REACT HOOK: useZoneHeatmap
 * Binds directly against our Backend's Persistent Syncer outputs.
 * Provides live WebSocket/Listen updates seamlessly integrated natively into the UI component tree.
 */
export const useZoneHeatmap = (zoneId: string) => {
  const [heatmap, setHeatmap] = useState<Heatmap | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    if (!zoneId) return;

    // Attaching dynamic real-time listener mapped efficiently avoiding explicit manual refreshing
    const docRef = doc(db, 'heatmaps', zoneId);
    const unsubscribe = onSnapshot(docRef, (docSnap) => {
      if (docSnap.exists()) {
        const data = docSnap.data();
        setHeatmap({
          zone_id: docSnap.id,
          density_level: data.DensityLevel,
          timestamp: data.Timestamp
        });
      }
      setLoading(false);
    }, (error) => {
      console.error("[SDK Error] Firestore native listener failed:", error);
      setLoading(false);
    });

    return () => unsubscribe();
  }, [zoneId]);

  return { heatmap, loading };
};

/**
 * CLIENT SDK CLASS: StadiumAPI
 * Strictly implements standard JSON:API boundaries bridging our Global HTTP(S) Load Balancer natively.
 */
export class StadiumAPI {
  private baseUrl = "https://api.stadium-experience.com/v1";

  // Extracts IAP context mapped token ensuring secure Zero-Trust architecture natively.
  constructor(private iapToken: string) {}

  async fetchWaitTime(amenityId: string): Promise<WaitTime> {
    const res = await fetch(`${this.baseUrl}/stalls/wait-times?amenity=${amenityId}`, {
      headers: {
        'x-goog-iap-jwt-assertion': this.iapToken,
      }
    });

    if (!res.ok) throw new Error("SDK WaitTime query execution critically failed");
    const json = await res.json();
    
    // Unwraps JSON:API structured attributes mapping straight strictly back to TypeScript interfaces
    return json.data.attributes as WaitTime;
  }
}
