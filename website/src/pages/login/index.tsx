import styles from './index.less';
import {Button, Typography} from "antd";
import ProForm, {ProFormInstance, ProFormText} from '@ant-design/pro-form';
import {useRef, useState} from "react";
import {POST} from "@/global";
import {KeyOutlined, UserOutlined, LockOutlined} from "@ant-design/icons";

const {Text, Title} = Typography;

const title = eval("window.config.title")
const desc = eval("window.config.desc")

export default function IndexPage() {
  const [btnLoading, setBtnLoading] = useState(false);
  const formRef = useRef<ProFormInstance>();

  const finish = (values: any) => {
    setBtnLoading(true);
    POST("/api/login", values, (res: any, status: any) => {
      setBtnLoading(false);
      if (status == 200) {
        window.location.href = "/";
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
                      onClick={() => formRef.current?.submit()}>登录</Button>
            ,
          }}
          onFinish={async (values) => finish(values)}
        >
          <ProForm.Group>
            <Title level={2}>{title}</Title>
          </ProForm.Group>

          <div className={styles.desc}>
            <Text type="secondary">{desc}</Text>
          </div>


          <ProForm.Group>
            <ProFormText
              width={"md"}
              tooltip={"最少3位字符，最大32位字符"}
              placeholder={"请输入登录账号"}
              fieldProps={{
                size: 'large',
                prefix: <UserOutlined/>,
              }}
              label="登录账号"
              rules={[
                {required: true, message: '请输入登录账号!'},
                {min: 3, message: '登录账号太短了!'},
                {max: 32, message: '登录账号太长了!'}
              ]}
              name="username"
            />
          </ProForm.Group>

          <ProForm.Group>
            <ProFormText.Password
              width={"md"}
              placeholder={"请输入登录密码"}
              tooltip={"最少6个字符，最大32个字符"}
              fieldProps={{
                size: 'large',
                prefix: <LockOutlined />,
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
        </ProForm>
      </div>
    </div>
  );
}
