import {PageContainer} from "@ant-design/pro-layout";
import {Alert, Card, Col, Row, Skeleton, Statistic} from "antd";
import {useEffect, useState} from "react";
import {GET} from "@/global";
import {Chart} from "@antv/g2";

export default function IndexPage(props: any) {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState({
    taskCount: 0,
    runCount: 0,
    schedulerCount: 0,
  })
  useEffect(() => {
    setLoading(true);
    GET(`/api/u/dashboard`, function (res: any, status: any) {
      setLoading(false);
      if (status == 200) {
        setData(res.data)
        if (res.data.chart.length > 0) {
          chart(res.data.chart);
        }
      }
    })
  }, []);

  const chart = (data: any) => {
    var data = data.map((item: any) => {
      if (item.code == 200) {
        item.code = "成功"
      } else {
        item.code = "失败"
      }
      return {
        ...item,
        code: item.code + ""
      }
    })
    const chart = new Chart({
      container: 'container',
      autoFit: true,
      height: 500,
    });

    chart.data(data);
    chart.scale({
      date: {
        range: [0, 1],
      },
      count: {
        nice: true,
      },
    });

    chart.tooltip({
      showCrosshairs: true,
      shared: true,
    });

    chart.axis('count', {
      label: {
        formatter: (val) => {
          return val
        },
      },
    });

    chart
      .line()
      .position('date*count')
      .color('code', (v) => {
        if (v == "失败") {
          return 'red';
        }
        return 'green';
      })
      .shape('smooth');

    chart
      .point()
      .position('date*count')
      .color('code', (v) => {
        if (v == "失败") {
          return 'red';
        }
        return 'green';
      })
      .shape('circle');

    chart.render();
  }

  return (
    <PageContainer
      key={"xxz"}
      // content="欢迎使用 ProLayout 组件"
      // tabList={[
      //   {
      //     tab: '基本信息',
      //     key: 'base',
      //   },
      //   {
      //     tab: '详细信息',
      //     key: 'info',
      //   },
      // ]}
      extra={[<a key="metrics" href={"/metrics"} target="_blank">Metrics</a>]}
      footer={[]}
    >
      <Skeleton loading={loading} active={true} paragraph={{rows: 18}}>
        <div className="site-statistic-demo-card">
          <Row gutter={16}>
            <Col span={12}>
              <Card>
                <Statistic
                  title="任务数量 / 运行数量"
                  value={data.taskCount + " / " + data.runCount}
                  precision={0}
                  // valueStyle={{color: '#3f8600'}}
                  // prefix={<ArrowUpOutlined/>}
                  // suffix="%"
                />
              </Card>
            </Col>
            <Col span={12}>
              <Card>
                <Statistic
                  title="累计调度次数"
                  value={data.schedulerCount}
                  precision={0}
                  // valueStyle={{color: '#cf1322'}}
                />
              </Card>
            </Col>
          </Row>

          {data.schedulerCount > 0 ? <Row gutter={16} style={{marginTop: 40}}>
            <Col span={24}>
              <div id="container"/>
            </Col>
          </Row> : <Alert message={"暂无调度数据"} type={"warning"} style={{marginTop: 20}}></Alert>}

        </div>
      </Skeleton>
    </PageContainer>
  );
}
