import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface MyQuery extends DataQuery {
  queryText?: string;
  serverUrl: string;
  withStreaming: boolean;
}

export const defaultQuery: Partial<MyQuery> = {
  serverUrl: "127.0.0.1:6570",
  withStreaming: false,
};

/**
 * These are options configured for each DataSource instance.
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  path?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  apiKey?: string;
}
