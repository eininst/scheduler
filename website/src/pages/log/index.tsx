import {PageContainer} from "@ant-design/pro-layout";
import type {ActionType, ProColumns} from '@ant-design/pro-components';
import {ProTable,} from '@ant-design/pro-components';
import {Tag, Typography} from 'antd';
import {GET} from "@/global";
import {useRef} from "react";
import ReactJson from "react-json-view";
const {Text, Link} = Typography;

export type TableListItem = {
  id: number;
  userId: number;
  userName: string;
  userRealName: string;
  userHead: string;
  taskId: number;
  taskName: string;
  taskGroup: string;
  taskUrl: string;
  taskObj: string;
  code: number;
  response: string;
  start_time: string;
  end_time: string;
  duration: number;
  createTime: number;
  obj: any;
};


const userRequest = async (params: any) => {
  var p = new Promise<any>((resolve) => {
    GET("/api/u/user", (res: any, status: any) => {
      if (status == 200) {
        var d = res.data.map((item: any) => {
          return {
            label: item.realName == "" ? item.userName : item.realName,
            value: item.id + "",
          }
        })
        resolve(d)
      } else {
        resolve([])
      }
    })
  })
  return await p;
};


export default function IndexPage(p:any) {
  const actionRef = useRef<ActionType>();
  const expandedRowRender = (r: any) => {
    var data = [{
      key: r.id,
      ...r
    }]
    return (
      <ProTable
        key={"exp" + r.id}
        columns={[
          {
            title: 'Request', dataIndex: 'req', key: 'req' + r.id, render: (_, record) => {
              var s: any = {
                "URL": record.taskUrl,
                "Method": record.obj.method,
                "Timeout": r.obj.timeout + "s",
                "MaxRetries": r.obj.maxRetries,
              }
              if (r.obj.contentType) {
                s["Content-Type"] = r.obj.contentType
              }
              if (r.obj.body) {
                s["Body"] = r.obj.contentType
              }

              return <ReactJson key={"j1" + r.id} src={s} name={null} collapsed={false}/>
            }
          },
          {
            title: 'Response', dataIndex: 'resp', key: 'resp' + r.id, render: (_, record) => {
              var resp = {
                Response: r.response,
              }
              return <ReactJson key={"j2" + r.id} src={resp} name={null} collapsed={false}/>
            }
          },
        ]}
        headerTitle={false}
        search={false}
        options={false}
        dataSource={data}
        pagination={false}
      />
    );
  }


  const columns: ProColumns<TableListItem>[] = [
    {
      width: 40,
      dataIndex: 'index',
      valueType: 'indexBorder',
    },
    {
      title: '??????',
      width: 80,
      hideInTable: true,
      dataIndex: 'taskGroup',
    },
    {
      title: '??????',
      width: 150,
      dataIndex: 'taskName',
      render: (_, record) => {
        if (record.taskGroup != '') {
          return <a>{record.taskGroup}???{_}</a>
        }
        return <a>{_}</a>
      },
    },


    {
      title: '??????',
      width: 80,
      hideInTable: true,
      dataIndex: 'status',
      valueEnum: {
        ok: {text: '?????????', status: 'Success'},
        fail: {text: '?????????', status: 'Error'},
      },
    },

    {
      title: 'Cron',
      dataIndex: 'cron',
      width: 110,
      search: false,
      copyable: true,
      render: (_, record) => {
        return record.obj.cron;
      },
    },
    {
      title: 'Url',
      dataIndex: 'taskUrl',
      align: 'left',
      search: false,
      copyable: true,
      render: (dom, record) => {
        return (<div><Text type={"secondary"}>{record.obj.method}: </Text> {dom}</div>)
      },
    },
    {
      title: 'Code',
      width: 60,
      dataIndex: 'code',
      render: (_, record) => {
        if (record.code >= 200 && record.code < 300) {
          return <Tag color="success">{record.code}</Tag>
        } else {
          return <Tag color="error">{record.code}</Tag>
        }
      },
    },

    {
      title: "??????",
      width: 70,
      align: "left",
      key: 'duration',
      // valueType: 'dateRange',
      dataIndex: 'duration',
      render: (_, record) => {
        return record.duration + "ms"
      },
      sorter: (a, b) => a.duration - b.duration,
    },
    {
      title: '?????????',
      width: 80,
      dataIndex: 'userId',
      fieldProps: {
        showSearch: true,
      },
      request: userRequest,
      render: (_, record) => {
        if (record.userRealName != '') {
          return record.userRealName;
        }
        return record.userName;
      },
    },
    {
      title: "????????????",
      width: 200,
      key: 'start_time',
      // valueType: 'dateRange',
      dataIndex: 'start_time',
      search: false,
      sorter: (a, b) => a.createTime - b.createTime,
    },

  ];

  return (
    <PageContainer>
      <ProTable<TableListItem>
        columns={columns}
        actionRef={actionRef}
        request={async (params = {}, sort, filter) => {
          for (var k in sort) {
            params['sort'] = k
            params['dir'] = sort[k]
            break
          }
          return new Promise((resolve, reject) => {
            GET("/api/u/task/excute/page", params, (res: any, status: number) => {
              if (status == 200) {
                var r = res.data;
                resolve({
                  data: r.list.map((item: any) => {
                    delete item['userHead']
                    delete item['userMail']
                    return {
                      ...item,
                      obj: JSON.parse(item.taskObj),
                    }
                  }),
                  // success ????????? true???
                  // ?????? table ???????????????????????????????????????
                  success: true,
                  // ??????????????? data ???????????????????????????????????????
                  total: r.total,
                })
              } else {
                resolve({data: [], success: false, total: 0})
              }
            })
          })
        }}

        form={{
          // ??????????????? transform????????????????????????????????????????????????????????????
          syncToUrl: (values, type) => {
            if (type === 'get') {
              return {
                ...values,
                created_at: [values.startTime, values.endTime],
              };
            }
            return values;
          },
        }}
        expandable={{expandedRowRender}}
        rowKey="id"
        pagination={{
          showQuickJumper: true,
          defaultPageSize: 10,
          showSizeChanger: true,
          pageSizeOptions: [10, 20]
        }}
        search={{
          filterType: 'light',
        }}
        // dateFormatter="string"
        headerTitle="????????????"
        toolBarRender={() => []}
      />

    </PageContainer>
  );
}
