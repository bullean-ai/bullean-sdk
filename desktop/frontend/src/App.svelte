
<main style="--default-contextmenu: show;">
  <div class="result" id="result">{resultText}</div>
  <div class="input-box" id="input">
    <!--<input autocomplete="off" bind:value={name} class="input" id="name" type="text"/>-->
    <button class="btn" on:click={getPredictions}>Evaluate Strategy</button>
  </div>
  <div id="container"></div>
</main>

<style>

  #logo {
    display: block;
    width: 50%;
    height: 50%;
    margin: auto;
    padding: 10% 0 0;
    background-position: center;
    background-repeat: no-repeat;
    background-size: 100% 100%;
    background-origin: content-box;
  }

  .result {
    height: 20px;
    line-height: 20px;
    margin: 1.5rem auto;
  }

  .input-box .btn {
    width: 60px;
    height: 30px;
    line-height: 30px;
    border-radius: 3px;
    border: none;
    margin: 0 0 0 20px;
    padding: 0 8px;
    cursor: pointer;
  }

  .input-box .btn:hover {
    background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
    color: #333333;
  }

  .input-box .input {
    border: none;
    border-radius: 3px;
    outline: none;
    height: 30px;
    line-height: 30px;
    padding: 0 10px;
    background-color: rgba(240, 240, 240, 1);
    -webkit-font-smoothing: antialiased;
  }

  .input-box .input:hover {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

  .input-box .input:focus {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

</style>

<script lang="ts">
  import {InitCandles, GetPredictions} from '../wailsjs/go/main/App.js'
  import {EventsOnce,EventsOn,EventsOnMultiple} from "../wailsjs/runtime";
  import Highcharts from 'highcharts/highstock';
  import dayjs from 'dayjs';

  let resultText: string = "Please enter your name below 👇. Yes Sir"
  let candles = []
  let predictions = []
  let name: string
  let chart;

 /*
   EventsOn("candles.init",result => {
    let res = JSON.parse(result)
    candles = [...candles, res]
    console.log(1)
  })
  */

  EventsOnce("candles.done",(result) => {
    console.log(1)
    if (result == true) {
      InitCandles("XRPUSDT").then(()=>{
        candles=[...candles]
        drawChart()
      })
    }
  })

  EventsOnMultiple("candles.new",result => {
    let res = JSON.parse(result)
    candles = [...candles, res]
    if (res != undefined || res != null) {
      drawChart()
    }
  },100000000)

  function getPredictions(): void {
    GetPredictions("XRPUSDT")
  }

  EventsOnMultiple("candles.prediction",result => {
    predictions = JSON.parse(result)
    drawChart()
  },100000000)

  function drawChart() {
    try {

      var option;
      let volume = []
      let chartData = []
      let pred = []
      const upColor = '#ec0000';
      const upBorderColor = '#8A0000';
      const downColor = '#00da3c';
      const downBorderColor = '#008F28';
      for (let i = 0; i<candles.length; i++) {
        let candle = candles[i]
        const time = new Highcharts.Time({
          timezone: 'Europe/Istanbul'
        });
        let dateInfo =  new Date(candle["t"])
        const s = time.dateFormat('%Y-%m-%d %H:%M:%S',Date.UTC(dateInfo.getUTCFullYear(),
          dateInfo.getUTCMonth(),
          dateInfo.getUTCDate(),
          dateInfo.getUTCHours(),
          dateInfo.getUTCMinutes(),
          dateInfo.getUTCSeconds()
        ));

        //let time = dayjs(new Date(candle["t"]).getTime()).format("dddd, MMM D, hh:mm::ss.SSS A-") //Thursday, Jan 1 at 12:00:01.090 AM-12:00:01.099 AM
        let data = [
          s,
          candle["o"],
          candle["h"],
          candle["l"],
          candle["c"],
        ]
        chartData.push(data)
        volume.push([
          time,
          candle["v"]
        ])
      }

      for (let i = 0; i<predictions.length; i++) {
        let prediction = predictions[i]
        //let time = dayjs(new Date(prediction["time"]).getTime()).format("dddd, MMM D, hh:mm a")
        let dateInfo =  new Date(prediction["time"])
        const time = new Highcharts.Time({
          timezone: 'Europe/Istanbul'
        });
        const s = time.dateFormat('%Y-%m-%d %H:%M:%S',Date.UTC(dateInfo.getUTCFullYear(),
                dateInfo.getUTCMonth(),
                dateInfo.getUTCDate(),
                dateInfo.getUTCHours(),
                dateInfo.getUTCMinutes(),
                dateInfo.getUTCSeconds()
        ));
        let data = [
          s,
          prediction["prediction"]
        ]
        pred.push(data)
      }
      console.log(candles)
      console.log(predictions)

      // create the chart
      chart = Highcharts.stockChart('container', {
        yAxis: [{
          labels: {
            align: 'left'
          },
          height: '80%',
          resize: {
            enabled: true
          }
        }, {
          labels: {
            align: 'left'
          },
          top: '80%',
          height: '20%',
          offset: 0
        }],
        rangeSelector: {
          selected: 1
        },
        xAxis:{
          zoomEnabled:true,
        },
        accessibility:{
          enabled:true
        },
        tooltip: {
          headerShape: 'callout',
          borderWidth: 0,
          shadow: false,
          fixed: true
        },
        series: [
          {
            type: 'candlestick',
            name: 'XRPUSDT',
            data: chartData,
          }, {
            type: 'column',
            id: 'predictions',
            name: 'Predictions',
            data: pred,
            yAxis: 1
          }
        ],
        annotations:[
          {
            draggable: "xy",
          }
        ],
        navigator:{
          enabled:true
        },
        responsive: {
          rules: [{
            condition: {
              maxWidth: 800
            },
            chartOptions: {
              rangeSelector: {
                inputEnabled: false
              }
            }
          }]
        }
      });
      /*
      // create the chart
      chart = Highcharts.chart('container', {
        rangeSelector: {
          selected: 1
        },
        xAxis:{
          zoomEnabled:true,
        },
        yAxis:{
          zoomEnabled:true,
        },
        accessibility:{
          enabled:true
        },
        title: {
          text: "XRPUSDT"
        },
        annotations: {
          draggable: "xy"
        },
        chart:{
          zooming:{
            key:"ctrl",
            pinchType:"x",
            singleTouch:true,
            mouseWheel:true,
            type:"x"
          }
        },
        navigator:{
          enabled:true
        },
        series: [{
          type: 'candlestick',
          name: 'XRPUSDT',
          data: chartData,
          dataGrouping: {
            units: [
              [
                'week', // unit name
                [1] // allowed multiples
              ], [
                'month',
                [1, 2, 3, 4, 6]
              ]
            ]
          }
        }]
      });
*/
    }catch (err) {
      alert(err)
    }

  }
</script>
