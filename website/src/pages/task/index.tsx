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
            label: item.realName == "" ? item.name : item.realName,
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
        message.success("????????????")
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
          message.success("????????????")
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
        message.success("????????????")
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
        message.success("????????????")
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
        message.success("???????????????" + res.data.count + " ???")
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
        message.success("???????????????" + res.data.count + " ???")
        actionRef.current?.reload();
      }
    })
  }

  const run = (id: number) => {
    setBtnLoading(true)
    POST("/api/u/task/start/" + id, {}, (res: any, status: any) => {
      setBtnLoading(false);
      if (status == 200) {
        message.success("????????????")
        actionRef.current?.reload();
      }
    })
  }

  const stop = (id: number) => {
    setBtnLoading(true)
    POST("/api/u/task/stop/" + id, {}, (res: any, status: any) => {
      setBtnLoading(false);
      if (status == 200) {
        message.success("????????????")
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
            message.success("????????????")
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
            message.success("????????????")
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
      title: '??????',
      dataIndex: 'group',
      width: 110,
      hideInTable: true,
    },
    {
      title: '????????????',
      width: 150,
      dataIndex: 'name',
      render: (_, record) => {
        if (record.group != '') {
          return <a>{record.group}???{_}</a>
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
      title: '??????',
      width: 80,
      dataIndex: 'status',
      valueEnum: {
        // run: {text: '??????', status: 'Default'},
        stop: {text: '?????????', status: 'Default'},
        run: {text: '?????????', status: 'Processing'},
      },
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
      title: (
        <>
          ????????????
          <Tooltip placement="top" title="??????????????????">
            <QuestionCircleOutlined style={{marginInlineStart: 4}}/>
          </Tooltip>
        </>
      ),
      width: 154,
      key: 'createTime',
      dataIndex: 'createTime',
      search: false,
      sorter: (a, b) => a.createTime - b.createTime,
    },

    {
      title: '??????',
      width: 180,
      key: 'option',
      valueType: 'option',
      render: (dom, record) => [
        (record.status == "stop" ?
          <a key={"run" + record.id} onClick={() => run(record.id)}>??????</a>
          : <a key={"run" + record.id} onClick={() => stop(record.id)}>??????</a>),
        <a key={"edit" + record.id} onClick={() => edit(record)}>??????</a>,
        <a key={"log" + record.id} onClick={()=> history.push("/log?taskName="+record.name)}>??????</a>,
        <TableDropdown
          key={"drop" + record.id}
          menus={[
            {
              key: 'do' + record.id, name: (
                <Popconfirm
                  key={"del" + record.id}
                  title="?????????????????????????"
                  onConfirm={() => {
                    doOnece(record.id);
                  }}
                  okText="???"
                  cancelText="???"
                >
                  <a href="#">????????????</a>
                </Popconfirm>
              )
            },
            {
              key: 'copy' + record.id, name: '???????????????', onClick: () => {
                showChangeUser(record)
              }
            },
            {
              key: 'delete' + record.id, name: (
                <Popconfirm
                  key={"del" + record.id}
                  title="???????????????????"
                  onConfirm={() => {
                    del(record.id);
                  }}
                  okText="???"
                  cancelText="???"
                >
                  <a href="#">??????</a>
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
        tableAlertOptionRender={(r) => {
          return (
            <Space size={16}>
              <a onClick={() => showBatchChangeUser(r.selectedRowKeys)}>???????????????</a>
              <a onClick={() => batchRun(r.selectedRowKeys)}>????????????</a>
              <a onClick={() => batchStop(r.selectedRowKeys)}>????????????</a>
              <a onClick={() => batchDel(r.selectedRowKeys)}>????????????</a>
            </Space>
          );
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
        headerTitle="????????????"
        toolBarRender={() => [
          <DrawerForm<{
            name: string;
            company: string;
          }>
            open={showEdit}
            title={isEdit ? "????????????" : "??????????????????"}
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
                ????????????
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

            {readonly?<Alert message="?????????????????????????????????" type="error" style={{marginBottom: 20}}/>:null}
            <ProFormText
              hidden={true}
              name="id"
            />
            <ProForm.Group>
              <ProFormText
                readonly={readonly}
                name="name"
                width="md"
                label="????????????"
                tooltip="????????? 32 ???"
                placeholder="?????????????????????"
                rules={[
                  {required: true, message: '?????????????????????!'},
                  {min: 1, message: '?????????!'},
                  {max: 32, message: '?????????!'}
                ]}
              />
              <ProFormText width="md" name="group" label="????????????" placeholder="?????????????????????(?????????)" readonly={readonly}/>
            </ProForm.Group>

            <ProForm.Group>
              <ProFormTextArea
                readonly={readonly}
                fieldProps={{
                  rows: 2,
                }}
                rules={[{max: 256, message: '?????????!'}]}
                width="xl" label="????????????" name="desc"/>
            </ProForm.Group>

            <ProForm.Group>
              <div className={styles.cron}>
                <ProFormText
                  readonly={readonly}
                  width="xl"
                  name="cron"
                  tooltip="Cron?????????, eg: */1 * * * * *"
                  label="Cron"
                  placeholder="?????????Cron?????????"
                  rules={[
                    {required: true, message: '?????????Cron?????????!'},
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
                  {required: true, message: '?????????Method!'},
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
                           placeholder="?????????Url"
                           rules={[
                             {required: true, message: '?????????Url!'},
                           ]}
              />
            </ProForm.Group>

            <ProForm.Group>
              <ProFormDigit width="xs" name="timeout" tooltip='??????"???"?????????60s'
                            label="????????????" initialValue={10}
                            readonly={readonly}
                            max={60}
                            min={1}
                            rules={[
                              {required: true, message: '?????????????????????!'},
                            ]}
              />

              <ProFormDigit width="xs" name="maxRetries" label="????????????" tooltip='"0" ???????????????????????? "10"'
                            readonly={readonly}
                            max={10}
                            min={0}
                            initialValue={3} rules={[
                {required: true, message: '?????????????????????!'},
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
                          {required: true, message: '?????????Method!'},
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
                          rules={[{max: 256, message: '?????????!'}]}
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
        title="?????????????????????"
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
        <Alert message="????????????????????????????????????" type="info" style={{marginBottom: 20}}/>
        <ProFormSelect
          request={userRequest}
          fieldProps={{
            showSearch: true,
          }}
          rules={[
            {required: true, message: '?????????????????????!'},
          ]}
          width="xl"
          name="userId"
          label="??????????????????"
        />
      </ModalForm>

    </PageContainer>
  );
}
