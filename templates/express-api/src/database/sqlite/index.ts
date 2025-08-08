import sqlite3 from 'sqlite3';
import { open } from 'sqlite';

let db: any = null;

const getDB = async () => {
  if (!db) {
    db = await open({
      filename: process.env.DB_PATH || './database.sqlite',
      driver: sqlite3.Database
    });
  }
  return db;
};

export default getDB; 