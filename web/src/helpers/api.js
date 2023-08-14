import { showError } from './utils';
import axios from 'axios';

export const API = axios.create({
  baseURL: process.env.REACT_APP_SERVER ? process.env.REACT_APP_SERVER : '',
});

API.interceptors.request.use(
    config => {
        if (localStorage.getItem('token')) {
            config.headers['Authorization'] = "Bearer " + localStorage.getItem('token');
        }
        return config
    },
    error => {
        console.log(error) // for debug
        return Promise.reject(error)
    }
)

API.interceptors.response.use(
  (response) => response,
  (error) => {
    showError(error);
  }
);
