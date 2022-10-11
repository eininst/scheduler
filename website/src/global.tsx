import {extend} from 'umi-request';
import {message} from "antd";

const req = extend({
  // prefix: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json;charset=UTF-8',
  },
});

export const POST = (uri: string, values: any, done: any) => {
  req.post(uri, {
    data: values,
    responseType: "json",
    getResponse: true,
    errorHandler: function (err) {
      if (err.response.status == 400) {
        message.error(err.data.message);
      }
      return err;
    }
  }).then(res => {
    done(res, res.response.status);
  });
}

export const GET = (uri: string, params: any, done?: any) => {
  if (typeof params == "function") {
    done = params;
  }
  req.get(uri, {
    responseType: "json",
    params: params,
    getResponse: true,
    errorHandler: function (err) {
      if (err.response.status == 400) {
        message.error(err.data.message);
      }
      return err;
    }
  }).then(res => {
    done(res, res.response.status);
  });
}

export const PUT = (uri: string, values: any, done: any) => {
  req.put(uri, {
    data: values,
    responseType: "json",
    getResponse: true,
    errorHandler: function (err) {
      if (err.response.status == 400) {
        message.error(err.data.message);
      }
      return err;
    }
  }).then(res => {
    done(res, res.response.status);
  });
}

export const DELETE = (uri: string, done: any) => {
  req.delete(uri, {
    responseType: "json",
    getResponse: true,
    errorHandler: function (err) {
      if (err.response.status == 400) {
        message.error(err.data.message);
      }
      return err;
    }
  }).then(res => {
    done(res, res.response.status);
  });
}
