import {PageContainer} from "@ant-design/pro-layout";
import {PlusOutlined, QuestionCircleOutlined} from '@ant-design/icons';
import type {ActionType, ProColumns} from '@ant-design/pro-components';
import {
  DrawerForm,
  ModalForm,
  ProForm,
  ProFormDependency,
  ProFormDigit,
  ProFormSelect,
  ProFormText,
  ProFormTextArea,
  ProTable,
  TableDropdown,
} from '@ant-design/pro-components';
import styles from './index.less';
import {Alert, Button, Form, message, Popconfirm, Space, Tooltip} from 'antd';
import {DELETE, GET, POST, PUT} from "@/global";
import {useRef, useState} from "react";
import ReactJson from "react-json-view";

import {history} from 'umi';

export type TableListItem = {
  id: number;
  name: string;
  group: string;
  cron: string;
  url: string;
  method: string;
  contentType: string;
  body: string;
  timeout: string;
  maxRetries: string;
  desc: string;
  status: string;
  userRealName: string;
  userName: string;
  createTime: number;
};

var role = eval("window.role")

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


var taskIds: any = [];
var isEdit: boolean = false;

export default function IndexPage() {
  const actionRef = useRef<ActionType>();
  const [formRef] = Form.useForm<any>();
  const [btnLoading, setBtnLoading] = useState(false);
  const [showEdit, setShowEdit] = useState(false);

  const [changeUser, setChangeUser] = useState(false);

  const [readonly, setReadonly] = useState(false);


  const edit = (record: any) => {
    setShowEdit(true)
    isEdit = true;
    formRef.setFieldsValue({...record})

    if (record.status == "run") {
      setReadonly(true);
    }else{
      setReadonly(false);
    }
  }

  const showChangeUser = (record: any) => {
    setChangeUser(true);
    taskIds = [record.id]
  }

  const showBatchChangeUser = (ids: any) => {
    setChangeUser(true);
    taskIds = ids
  }

  const doOnece = (id:number) =>{
    setBtnLoading(true)
    POST("/api/u/task/do/"+id, {
    }, (res: any, status: any) => {
      setBtnLoading(false);
      if (status == 200) {
        message.success("执行成功")
      }
    })
  }

  const submitChangeUser = (values: any) => {
    var userId = values.userId;
    return new Promise((resolve) => {
      POST("/api/u/task/batch/change/user", {
        userId: parseInt(userId),
        taskIds: taskIds,
      }, (res: any, status: any) => {
        setBtnLoading(false);
        if (status == 200) {
          message.success("操作成功")
          resolve(true);
          actionRef.current?.reload();
        } else {
          resolve(false);
        }
      })
    });
  }

  const del = (id: number) => {
    setBtnLoading(true)
    DELETE("/api/u/task/del/" + id, (res: any, status: any) => {
      setBtnLoading(false);
      if (status == 200) {
        message.success("操作成功")
        actionRef.current?.reload();
      }
    })
  }

  const batchDel = (ids: any) => {
    setBtnLoading(true)
    POST("/api/u/task/batch/del", {
      taskIds: ids,
    }, (res: any, status: any) => {
      setBtnLoading(false);
      if (status == 200) {
        message.success("操作成功")
        actionRef.current?.reload();
      }
    })
  }

  const batchRun = (ids: any) => {
    setBtnLoading(true)
    POST("/api/u/task/batch/start", {
      taskIds: ids,
    }, (res: any, status: any) => {
      setBtnLoading(false);
      if (status == 200) {
        message.success("成功启动：" + res.data.count + " 个")
        actionRef.current?.reload();
      }
    })
  }

  const batchStop = (ids: any) => {
    setBtnLoading(true)
    POST("/api/u/task/batch/stop", {
      taskIds: ids,
    }, (res: any, status: any) => {
      setBtnLoading(false);
      if (status == 200) {
        message.success("成功停止：" + res.data.count + " 个")
        actionRef.current?.reload();
      }
    })
  }

  const run = (id: number) => {
    setBtnLoading(true)
    POST("/api/u/task/start/" + id, {}, (res: any, status: any) => {
      setBtnLoading(false);
      if (status == 200) {
        message.success("启动成功")
        actionRef.current?.reload();
      }
    })
  }

  const stop = (id: number) => {
    setBtnLoading(true)
    POST("/api/u/task/stop/" + id, {}, (res: any, status: any) => {
      setBtnLoading(false);
      if (status == 200) {
        message.success("停止成功")
        actionRef.current?.reload();
      }
    })
  }

  const finish = (values: any) => {
    setBtnLoading(true);
    if (isEdit) {
      return new Promise((resolve) => {
        PUT("/api/u/task/update", values, (res: any, status: any) => {
          setBtnLoading(false);
          if (status == 200) {
            message.success("操作成功")
            resolve(true);
            actionRef.current?.reload();
          } else {
            resolve(false);
          }
        })
      });
    } else {
      return new Promise((resolve) => {
        POST("/api/u/task/add", values, (res: any, status: any) => {
          setBtnLoading(false);
          if (status == 200) {
            message.success("操作成功")
            resolve(true);
            actionRef.current?.reload();
          } else {
            resolve(false);
          }
        })
      });
    }
  }

  const expandedRowRender = (r: any) => {
    var json = {
      "Method": r.method,
      "Content-Type": r.contentType,
      "Timeout": r.timeout + "s",
      "MaxRetries": r.maxRetries,
      "Desc": r.desc,
    }
    return (<ReactJson src={r} name={null} collapsed={false}/>)
  }


  const columns: ProColumns<TableListItem>[] = [
    {
      width: 40,
      dataIndex: 'index',
      valueType: 'indexBorder',
    },
    {
      title: '分组',
      dataIndex: 'group',
      width: 110,
      hideInTable: true,
    },
    {
      title: '任务名称',
      width: 150,
      dataIndex: 'name',
      render: (_, record) => {
        if (record.group != '') {
          return <a>{record.group}：{_}</a>
        }
        return <a>{_}</a>
      },
    },
    {
      title: 'Cron',
      dataIndex: 'cron',
      width: 110,
      search: false,
      copyable: false,
    },
    {
      title: 'Url',
      dataIndex: 'url',
      align: 'left',
      search: false,
      copyable: true,
    },
    {
      title: '状态',
      width: 80,
      dataIndex: 'status',
      valueEnum: {
        // run: {text: '全部', status: 'Default'},
        stop: {text: '已停止', status: 'Default'},
        run: {text: '运行中', status: 'Processing'},
      },
    },
    {
      title: '创建者',
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
      title: (
        <>
          创建时间
          <Tooltip placement="top" title="这是一段描述">
            <QuestionCircleOutlined style={{marginInlineStart: 4}}/>
          </Tooltip>
        </>
      ),
      width: 170,
      key: 'createTime',
      // valueType: 'dateRange',
      dataIndex: 'createTime',
      search: false,
      sorter: (a, b) => a.createTime - b.createTime,
    },

    {
      title: '操作',
      width: 180,
      key: 'option',
      valueType: 'option',
      render: (dom, record) => [
        (record.status == "stop" ?
          <a key={"run" + record.id} onClick={() => run(record.id)}>启动</a>
          : <a key={"run" + record.id} onClick={() => stop(record.id)}>停止</a>),
        <a key={"edit" + record.id} onClick={() => edit(record)}>修改</a>,
        <a key={"log" + record.id} onClick={()=> history.push("/log?taskName="+record.name)}>日志</a>,
        <TableDropdown
          key={"drop" + record.id}
          menus={[
            {
              key: 'do' + record.id, name: (
                <Popconfirm
                  key={"del" + record.id}
                  title="确认要执行一次吗?"
                  onConfirm={() => {
                    doOnece(record.id);
                  }}
                  okText="是"
                  cancelText="否"
                >
                  <a href="#">执行一次</a>
                </Popconfirm>
              )
            },
            {
              key: 'copy' + record.id, name: '变更创建者', onClick: () => {
                showChangeUser(record)
              }
            },
            {
              key: 'delete' + record.id, name: (
                <Popconfirm
                  key={"del" + record.id}
                  title="确认要删除吗?"
                  onConfirm={() => {
                    del(record.id);
                  }}
                  okText="是"
                  cancelText="否"
                >
                  <a href="#">删除</a>
                </Popconfirm>
              )
            },
          ]}
        />,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<TableListItem>
        columns={columns}
        loading={btnLoading}
        actionRef={actionRef}
        request={async (params = {}, sort, filter) => {
          for (var k in sort) {
            params['sort'] = k
            params['dir'] = sort[k]
            break
          }
          return new Promise((resolve, reject) => {
            setBtnLoading(true)
            GET("/api/u/task/page", params, (res: any, status: number) => {
              setBtnLoading(false)
              if (status == 200) {
                var r = res.data;
                resolve({
                  data: r.list.map((item: any) => {
                    delete item['userHead']
                    delete item['userMail']
                    return item
                  }),
                  // success 请返回 true，
                  // 不然 table 会停止解析数据，即使有数据
                  success: true,
                  // 不传会使用 data 的长度，如果是分页一定要传
                  total: r.total,
                })
              } else {
                resolve({data: [], success: false, total: 0})
              }
            })
          })
        }}
        tableAlertOptionRender={(r) => {
          return (
            <Space size={16}>
              <a onClick={() => showBatchChangeUser(r.selectedRowKeys)}>变更创建者</a>
              <a onClick={() => batchRun(r.selectedRowKeys)}>批量启动</a>
              <a onClick={() => batchStop(r.selectedRowKeys)}>批量停止</a>
              <a onClick={() => batchDel(r.selectedRowKeys)}>批量删除</a>
            </Space>
          );
        }}
        form={{
          // 由于配置了 transform，提交的参与与定义的不同这里需要转化一下
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
        rowSelection={role=="admin"?{}:false}
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
        headerTitle="任务列表"
        toolBarRender={() => [
          <DrawerForm<{
            name: string;
            company: string;
          }>
            open={showEdit}
            title={isEdit ? "修改任务" : "创建一个任务"}
            form={formRef}
            onOpenChange={(b) => {
              setShowEdit(b)
              if (!b) {
                isEdit = false;
                setReadonly(false);
              }
            }}
            trigger={
              <Button type="primary" loading={btnLoading}>
                <PlusOutlined/>
                创建任务
              </Button>
            }
            autoFocusFirstInput
            drawerProps={{
              destroyOnClose: true,
            }}
            submitTimeout={2000}
            onFinish={async (values): Promise<any> => {
              var s = await finish(values);
              return s;
            }}
          >

            {readonly?<Alert message="运行中的任务不支持修改" type="error" style={{marginBottom: 20}}/>:null}
            <ProFormText
              hidden={true}
              name="id"
            />
            <ProForm.Group>
              <ProFormText
                readonly={readonly}
                name="name"
                width="md"
                label="任务名称"
                tooltip="最长为 32 位"
                placeholder="请输入任务名称"
                rules={[
                  {required: true, message: '请输入任务名称!'},
                  {min: 1, message: '太短了!'},
                  {max: 32, message: '太长了!'}
                ]}
              />
              <ProFormText width="md" name="group" label="任务分组" placeholder="请输入分组名称(非必填)" readonly={readonly}/>
            </ProForm.Group>

            <ProForm.Group>
              <ProFormTextArea
                readonly={readonly}
                fieldProps={{
                  rows: 2,
                }}
                rules={[{max: 256, message: '太长了!'}]}
                width="xl" label="任务描述" name="desc"/>
            </ProForm.Group>

            <ProForm.Group>
              <div className={styles.cron}>
                <ProFormText
                  readonly={readonly}
                  width="xl"
                  name="cron"
                  tooltip="Cron表达式, eg: */1 * * * * *"
                  label="Cron"
                  placeholder="请输入Cron表达式"
                  rules={[
                    {required: true, message: '请输入Cron表达式!'},
                  ]}
                />
              </div>
            </ProForm.Group>

            <ProForm.Group>
              <ProFormSelect
                initialValue={"GET"}
                key={"www"}
                readonly={readonly}
                options={[
                  {
                    key: "GET",
                    value: 'GET',
                    label: 'GET',
                  },
                  {
                    key: "POST",
                    value: 'POST',
                    label: 'POST',
                  },
                  {
                    key: "PUT",
                    value: 'PUT',
                    label: 'PUT',
                  },
                  {
                    key: "DELETE",
                    value: 'DELETE',
                    label: 'DELETE',
                  },
                ]}
                rules={[
                  {required: true, message: '请选择Method!'},
                ]}
                width="xs"
                name="method"
                label="Method"
              />

              <ProFormText width="xl"
                           readonly={readonly}
                           name="url"
                           tooltip="eg: http://xxx.com/api/v1/test"
                           label="Url"
                           placeholder="请输入Url"
                           rules={[
                             {required: true, message: '请输入Url!'},
                           ]}
              />
            </ProForm.Group>

            <ProForm.Group>
              <ProFormDigit width="xs" name="timeout" tooltip='单位"秒"，最大60s'
                            label="超时时间" initialValue={10}
                            readonly={readonly}
                            max={60}
                            min={1}
                            rules={[
                              {required: true, message: '请输入超时时间!'},
                            ]}
              />

              <ProFormDigit width="xs" name="maxRetries" label="重试次数" tooltip='"0" 代表不重试，最大 "10"'
                            readonly={readonly}
                            max={10}
                            min={0}
                            initialValue={3} rules={[
                {required: true, message: '请输入重试次数!'},
              ]}/>
            </ProForm.Group>


            <ProFormDependency name={['method']}>
              {({method}) => {
                if (method != 'GET') {
                  return ([<ProForm.Group key={"k1"}>
                      <ProFormSelect
                        initialValue={"application/json"}
                        options={[
                          {
                            value: 'application/json',
                            label: 'application/json',
                          },
                          {
                            value: 'application/x-www-form-urlencoded',
                            label: 'application/x-www-form-urlencoded',
                          },
                        ]}
                        rules={[
                          {required: true, message: '请选择Method!'},
                        ]}
                        width="md"
                        readonly={readonly}
                        name="contentType"
                        label="Content-Type"
                      />
                    </ProForm.Group>,

                      <ProForm.Group key={"k2"}>
                        <ProFormTextArea
                          readonly={readonly}
                          fieldProps={{
                            rows: 3,
                          }}
                          rules={[{max: 256, message: '太长了!'}]}
                          width="xl" label="Body" name="body"/>
                      </ProForm.Group>]
                  )
                }
              }}
            </ProFormDependency>

            {/*<ReactJson src={my_json_object} name={null} collapsed={false} />*/}
          </DrawerForm>,
        ]}
      />


      <ModalForm<{
        name: string;
        company: string;
      }>
        width={600}
        title="变更任务创建者"
        open={changeUser}
        onOpenChange={(b) => {
          setChangeUser(b);
          if (!b) {
            taskIds = [];
          }
        }}
        autoFocusFirstInput
        modalProps={{
          destroyOnClose: true,
          onCancel: () => console.log('run'),
        }}
        submitTimeout={2000}
        onFinish={async (values: any): Promise<any> => {
          return await submitChangeUser(values);
        }}
      >
        <Alert message="请选择一个用户，进行变更" type="info" style={{marginBottom: 20}}/>
        <ProFormSelect
          request={userRequest}
          fieldProps={{
            showSearch: true,
          }}
          rules={[
            {required: true, message: '请选择一个用户!'},
          ]}
          width="xl"
          name="userId"
          label="选择一个用户"
        />
      </ModalForm>

    </PageContainer>
  );
}
