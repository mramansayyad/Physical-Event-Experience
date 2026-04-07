import os
import xgboost as xgb
import pandas as pd
from google.cloud import bigquery
from google.cloud import storage

def extract_historical_features():
    client = bigquery.Client()
    
    # Extract structural density aggregates over a rolling 30-minute block explicitly predicting future heat
    query = """
    SELECT 
        zone_id,
        EXTRACT(HOUR FROM timestamp) as hour,
        EXTRACT(MINUTE FROM timestamp) as minute,
        COUNT(DISTINCT device_id) as current_density,
        LEAD(COUNT(DISTINCT device_id), 6) OVER (
            PARTITION BY match_id, zone_id 
            ORDER BY TIMESTAMP_TRUNC(timestamp, MINUTE)
        ) as density_in_30_mins
    FROM `stadium-experience-loc.analytics.telemetry_stream`
    GROUP BY match_id, zone_id, timestamp
    """
    
    df = client.query(query).to_dataframe()
    return df.dropna()

def train_xgboost_hotspot_model():
    print("Initiating Vertex AI Training Script via XGBoost natively...")
    dataset = extract_historical_features()
    
    # Feature matrix extracting current dynamics
    X = dataset[['current_density', 'hour', 'minute']]
    y = dataset['density_in_30_mins']
    
    # Execute native XGBoost regression bounding highly specialized predictive accuracy natively
    model = xgb.XGBRegressor(objective='reg:squarederror', n_estimators=150, max_depth=6)
    model.fit(X, y)
    
    # Save standard framework bounds securely mapping GCS uploads for Vertex Deployment natively
    model.save_model("stadium_congestion_model.json")
    print("Optimization Complete. Model serialized for continuous API serving.")

    # Upload to Cloud Storage natively initializing Vertex AI endpoints
    storage_client = storage.Client()
    bucket = storage_client.bucket("stadium-vertex-models-secure")
    blob = bucket.blob("latest/stadium_congestion_model.json")
    blob.upload_from_filename("stadium_congestion_model.json")

if __name__ == "__main__":
    train_xgboost_hotspot_model()
