import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ChartConfiguration, ChartOptions } from 'chart.js';
import { Subscription } from 'rxjs';
import { ApiService, Report, Stat } from '../service/api.service';

@Component({
  selector: 'app-details',
  templateUrl: './details.component.html',
  styleUrls: ['./details.component.scss']
})
export class DetailsComponent implements OnInit, OnDestroy {

  private subscriptions: Subscription[] = []
  unsafeStats: Stat[] = []
  query: string | null = null
  hasUnsafeStats: boolean = false
  parsedReports: ParsedReport[] = []
  title = 'ng2-charts-demo';
  hasReports = false
  graphicReports: GraphicReport[] = []

  public lineChartData: ChartConfiguration<'line'>['data'] = {
    labels: [],
    datasets: [
      {
        data: [],
        label: 'Series A',
        fill: true,
        tension: 0.5,
        borderColor: 'black',
        backgroundColor: 'rgba(255,0,0,0.3)'
      }
    ]
  };
  public lineChartOptions: ChartOptions<'line'> = {
    responsive: false,
    scales: {
      x: {
        ticks: {
          display: false
        }
      }
    }
  };
  public lineChartLegend = false;


  public constructor(private api: ApiService, private route: ActivatedRoute) {

  }

  ngOnDestroy(): void {
    this.subscriptions.forEach(x => x.unsubscribe)
  }

  ngOnInit(): void {
    let routeSubscription = this.route.paramMap.subscribe(paramMap => {
      this.query = paramMap.get('query')
      this.getDetails();
    })

    this.subscriptions.push(routeSubscription)
  }

  getDetails(): void {
    if (this.query === null) {
      return;
    }
    let unsafeDetailsSubscription = this.api.getUnsafeDetails(this.query).subscribe( (data: Stat[]) => {
      this.unsafeStats = data
      this.mapUnsafeStats()
    })
    this.subscriptions.push(unsafeDetailsSubscription)
    let reportsSubscription = this.api.getReports(this.query).subscribe( (data: Report[]) => {
      let parsedReports = data.map((x): ParsedReport => {
        x.data = x.data.replaceAll('\'', '')
        let rows = x.data.split('\n')
        let topics = rows.map(x => x.split("\t")).map((x): Topic => {
          return {name: x[0], value: Number(x[1])}
        }).filter(x => {
            const tooSmall = x.name.length <= 1
            const isInQuery = this.query?.toLocaleLowerCase().indexOf(x.name.toLocaleLowerCase())
            const lengthIsSimilar = Math.abs((this.query?.length ?? 0) - x.name.length) <= 1

            const isRedundant = isInQuery && lengthIsSimilar;

            return !(tooSmall || isRedundant)
        })
        topics = topics.sort((l, r) => r.value - l.value)
        if (topics.length > 30) {
          topics = topics.slice(0, 30)
        }
        return {id: x.id, query: x.query, time: new Date(x.timestamp/1000000), topics: topics}
      });
      
      this.parsedReports = parsedReports.sort((l, r) => r.time.getMilliseconds() - l.time.getMilliseconds())
      let graphicReports = this.parsedReports.map((x): GraphicReport => {
        return {origin: x, chartData: this.constructDatasetForReport(x)}
      })
      this.graphicReports = graphicReports;
      this.hasReports = true;
    })
    this.subscriptions.push(reportsSubscription)
  }

  generateReport() {
    if (this.query === null) {
      return;
    }
    let sub = this.api.triggerReport(this.query).subscribe(() => {console.log('Successfully scheduled report')})
    this.subscriptions.push(sub)
  }

  private mapUnsafeStats() {
    let sortedStats = this.unsafeStats.sort((left: Stat, right: Stat) => {
      return left.timestamp - right.timestamp;
    })
    if (sortedStats.length > 30) {
      sortedStats = sortedStats.slice(sortedStats.length - 30)
    }
    let values = sortedStats.map(x => x.value);
    let labels = sortedStats.map(x => {
      let t = new Date(1970, 0, 1)
      t.setSeconds(x.timestamp)
      return t
    })
    this.lineChartData.datasets[0].data = values;
    this.lineChartData.labels = labels
    this.hasUnsafeStats = true;
  }


  private constructDatasetForReport(report: ParsedReport): ChartConfiguration<'bar'>['data'] {
    return {
      labels: report.topics.map(x => x.name),
      datasets: [
        {
          data: report.topics.map(x => x.value),
          label: 'Series A',
          borderColor: 'black',
          backgroundColor: 'rgba(255,0,0,0.3)'
        }
      ]
    }
  }



  public reportChartOptions: ChartOptions<'bar'> = {
    responsive: false,
    scales: {
      x: {
        ticks: {
          display: false
        }
      }
    }
  };
}

export interface Topic {
  name: string,
  value: number
}

export interface ParsedReport {
  id: string
  time: Date,
  query: string,
  topics: Topic[]
}

export interface GraphicReport {
  origin: ParsedReport,
  chartData: ChartConfiguration<'bar'>['data']
}