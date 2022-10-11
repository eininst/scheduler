import styles from './index.less';
import {Button, Divider, message, Typography} from "antd";
import ProForm, {ProFormInstance, ProFormRadio, ProFormText} from '@ant-design/pro-form';
import {useRef, useState} from "react";
import {POST} from "@/global";

const {Text, Title} = Typography;

export default function IndexPage() {
  const [btnLoading, setBtnLoading] = useState(false);
  const formRef = useRef<ProFormInstance>();

  const finish = (values: any) => {
    if (values.password != values.password2) {
      message.error("两次输入的密码不一致");
      return
    }

    setBtnLoading(true);
    POST("/api/init", values, (res: any, status: any) => {
      setBtnLoading(false);
      if (status == 200) {
        window.location.href = "/"
      }
    })
  }

  return (
    <div className={styles.wrapper}>

      <div className={styles.spec}>
        <ProForm
          formRef={formRef}
          submitter={{
            render: (_, dom) =>
              <Button style={{width: "100%", marginTop: 10}} size={"large"} loading={btnLoading} type={"primary"}
                      onClick={() => formRef.current?.submit()}>Go</Button>
            ,
          }}
          onFinish={async (values) => finish(values)}
        >
          <ProForm.Group>
            <Title level={2}>初始化配置</Title>
          </ProForm.Group>

          <Divider orientation="left" plain style={{marginBottom: 20}}>
            初始化一个管理员账号
          </Divider>

          <ProForm.Group>
            <ProFormText
              width={"md"}
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

          <ProForm.Group>
            <ProFormText.Password
              width={"md"}
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
              width={"md"}
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

          <ProForm.Group>
            <ProFormRadio.Group
              label="账号权限"
              name="role"
              initialValue="admin"
              rules={[
                {required: true, message: '请选择账号权限!'},
              ]}
              options={[{
                value: 'admin',
                label: '管理员',
              }, {
                value: 'normal',
                label: '普通用户',
                disabled: true,
              }]}
            />
          </ProForm.Group>

          <ProForm.Group>
            <ProFormText
              width={"md"}
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
              width={"md"}
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
        </ProForm>
      </div>
    </div>
  );
}
