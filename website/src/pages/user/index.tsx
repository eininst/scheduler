import {PageContainer} from "@ant-design/pro-layout";
import {ActionType, DrawerForm, ModalForm, ProColumns, ProTable,} from '@ant-design/pro-components';
import {Alert, Button, Form, message, Popconfirm, Tooltip, Typography} from 'antd';
import {DELETE, GET, POST, PUT} from "@/global";
import {PlusOutlined, QuestionCircleOutlined} from "@ant-design/icons";
import {useRef, useState} from "react";
import ProForm, {ProFormRadio, ProFormText} from "@ant-design/pro-form";

const {Text, Title} = Typography;

export type TableListItem = {
  id: number;
  name: string;
  password: string;
  realName: string;
  role: string;
  mail: string;
  status: string;
  createTime: number;
};

var isEdit: boolean = false;
var id: number = 0;

export default function IndexPage() {
  const actionRef = useRef<ActionType>();
  const [formRef] = Form.useForm<any>();
  const [btnLoading, setBtnLoading] = useState(false);

  const [showEdit, setShowEdit] = useState(false);

  const [showReset, setShowRest] = useState(false);

  const [isAdd, setIsAdd] = useState(true);

  const enbale = (id: number) => {
    setBtnLoading(true);
    POST("/api/u/user/enable/" + id, {}, (res: any, status: any) => {
      setBtnLoading(false);
      if (status == 200) {
        message.success("操作成功")
        actionRef.current?.reload();
      }
    })
  }
  const disable = (id: number) => {
    setBtnLoading(true);
    POST("/api/u/user/disable/" + id, {}, (res: any, status: any) => {
      setBtnLoading(false);
      if (status == 200) {
        message.success("操作成功")
        actionRef.current?.reload();
      }
    })
  }

  const edit = (record: any) => {
    setShowEdit(true);
    setIsAdd(false);
    isEdit = true;
    formRef.setFieldsValue({...record, password2: record.password})
  }


  const showResetPassword = (record: any) => {
    setShowRest(true);
    id = record.id;
  }
  const resetPassword = (values: any) => {
    if (values.password2 != values.password) {
      message.error("两次输入的密码不一致")
      return
    }

    setBtnLoading(true);
    return new Promise((resolve) => {
      PUT("/api/u/user/password/reset", {
        id: id,
        password: values.password,
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

  const del = (id: any) => {
    setBtnLoading(true);
    DELETE("/api/u/user/del/" + id, (res: any, status: any) => {
      setBtnLoading(false);
      if (status == 200) {
        message.success("操作成功")
        actionRef.current?.reload();
      }
    })
  }

  const finish = (values: any) => {
    setBtnLoading(true);
    if (isEdit) {
      return new Promise((resolve) => {
        PUT("/api/u/user/update", values, (res: any, status: any) => {
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
      if (values.password2 != values.password) {
        message.error("两次输入的密码不一致")
        return
      }
      return new Promise((resolve) => {
        POST("/api/u/user/add", values, (res: any, status: any) => {
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

  const columns: ProColumns<TableListItem>[] = [
    {
      width: 40,
      dataIndex: 'index',
      valueType: 'indexBorder',
    },
    {
      title: '用户名',
      dataIndex: 'name',
      copyable: true,
      order:100,
      render: (_, record) => {
        return <a>{_}</a>
      },
    },
    {
      title: '权限',
      dataIndex: 'role',
      width: 110,
      copyable: false,
      valueEnum: {
        admin: {text: '管理员',},
        normal: {text: '普通用户',},
      },
    },
    {
      title: '真实姓名',
      dataIndex: 'realName',
      align: 'left',
      copyable: false,
      order:99,
    },
    {
      title: '报警邮箱',
      dataIndex: 'mail',
      align: 'left',
      search: false,
      copyable: false,
    },
    {
      title: '状态',
      width: 80,
      dataIndex: 'status',
      valueEnum: {
        ok: {text: '使用中', status: 'Success'},
        disabled: {text: '禁用中', status: 'Default'},
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
      width: 220,
      key: 'option',
      valueType: 'option',
      render: (dom, record) => [
        (record.status == "ok" ?
          <a key={"run" + record.id} onClick={() => disable(record.id)}>禁用</a>
          : <a key={"run" + record.id} onClick={() => enbale(record.id)}>启用</a>),
        <a key={"edit" + record.id} onClick={() => edit(record)}>修改</a>,
        <a key={"edit" + record.id} onClick={() => showResetPassword(record)}>重置密码</a>,
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
        </Popconfirm>,
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
            GET("/api/u/user", params, (res: any, status: number) => {
              setBtnLoading(false)
              if (status == 200) {
                var r = res.data;
                resolve({
                  data: r,
                  // success 请返回 true，
                  // 不然 table 会停止解析数据，即使有数据
                  success: true,
                  // 不传会使用 data 的长度，如果是分页一定要传
                  total: r.length,
                })
              } else {
                resolve({data: [], success: false, total: 0})
              }
            })
          })
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
        // expandable={{expandedRowRender}}
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
            title={isEdit ? "修改用户" : "创建一个用户"}
            form={formRef}
            onOpenChange={(b) => {
              setShowEdit(b)
              if (!b) {
                isEdit = false;
                setIsAdd(true)
              }
            }}
            trigger={
              <Button type="primary" loading={btnLoading}>
                <PlusOutlined/>
                创建一个用户
              </Button>
            }
            autoFocusFirstInput
            drawerProps={{
              destroyOnClose: true,
            }}
            submitTimeout={2000}
            onFinish={async (values): Promise<any> => {
              return await finish(values);
            }}
          >
            <ProFormText
              hidden={true}
              name="id"
            />

            <ProForm.Group>
              <ProFormText
                width={"lg"}
                tooltip={"最少3位字符，最大32位字符"}
                fieldProps={{
                  size: 'large',
                  // prefix: "% ",
                }}
                label="登录账号"
                rules={[
                  {required: true, message: '请输入登录账号!'},
                  {min: 3, message: '登录账号太短了!'},
                  {max: 32, message: '登录账号太长了!'}
                ]}
                name="name"
              />
            </ProForm.Group>

            {isAdd ? <div>
              <ProForm.Group>
                <ProFormText.Password
                  width={"lg"}
                  tooltip={"最少6个字符，最大32个字符"}
                  fieldProps={{
                    size: 'large',
                  }}
                  label="登录密码"
                  rules={[
                    {required: true, message: '请输入登录密码!'},
                    {min: 6, message: '密码太短了!'},
                    {max: 32, message: '密码太长了!'}
                  ]}
                  name="password"
                />
              </ProForm.Group>

              <ProForm.Group>
                <ProFormText.Password
                  width={"lg"}
                  tooltip={"再次输入登录密码"}
                  fieldProps={{
                    size: 'large',
                  }}
                  label="再次输入登录密码"
                  rules={[
                    {required: true, message: '请重复输入登录密码!'},
                    {min: 6, message: '密码太短了!'},
                    {max: 32, message: '密码太长了!'}
                  ]}
                  name="password2"
                />
              </ProForm.Group>
            </div> : null
            }

            <ProForm.Group>
              <ProFormRadio.Group
                label="账号权限"
                name="role"
                initialValue="normal"
                rules={[
                  {required: true, message: '请选择账号权限!'},
                ]}
                options={[{
                  value: 'admin',
                  label: '管理员',
                }, {
                  value: 'normal',
                  label: '普通用户',
                  // disabled: true,
                }]}
              />
            </ProForm.Group>

            <ProForm.Group>
              <ProFormText
                width={"lg"}
                tooltip={"最少3位字符，最大32位字符"}
                fieldProps={{
                  size: 'large',
                }}
                placeholder={"请输入真实姓名(非必填)"}
                label="真实姓名"
                rules={[
                  {min: 1, message: '太短了!'},
                  {max: 32, message: '太长了!'}
                ]}
                name="realName"
              />
            </ProForm.Group>

            <ProForm.Group>
              <ProFormText
                width={"lg"}
                tooltip={"任务出错时, 会统一抄送至管理员邮箱"}
                fieldProps={{
                  size: 'large',
                  // prefix: "% ",
                }}
                placeholder={"请输入报警邮箱(非必填)"}
                label="报警邮箱"
                rules={[
                  {type: 'email'},
                ]}
                name="mail"
              />
            </ProForm.Group>
          </DrawerForm>
        ]}
      />

      <ModalForm
        width={550}
        title="重置密码"
        open={showReset}
        onOpenChange={(b) => {
          setShowRest(b)
        }}
        autoFocusFirstInput
        modalProps={{
          destroyOnClose: true,
          onCancel: () => console.log('run'),
        }}
        submitTimeout={2000}
        onFinish={async (values): Promise<any> => {
          return await resetPassword(values);
        }}
      >
        <Alert message={"请输入新的密码，进行重置"} type="info" style={{marginBottom: 20}}></Alert>

        <ProForm.Group>
          <ProFormText.Password
            width={"xl"}
            tooltip={"最少6个字符，最大32个字符"}
            fieldProps={{
              size: 'large',
            }}
            label="登录密码"
            rules={[
              {required: true, message: '请输入登录密码!'},
              {min: 6, message: '密码太短了!'},
              {max: 32, message: '密码太长了!'}
            ]}
            name="password"
          />
        </ProForm.Group>

        <ProForm.Group>
          <ProFormText.Password
            width={"xl"}
            tooltip={"再次输入登录密码"}
            fieldProps={{
              size: 'large',
            }}
            label="再次输入登录密码"
            rules={[
              {required: true, message: '请重复输入登录密码!'},
              {min: 6, message: '密码太短了!'},
              {max: 32, message: '密码太长了!'}
            ]}
            name="password2"
          />
        </ProForm.Group>
      </ModalForm>

    </PageContainer>
  );
}
