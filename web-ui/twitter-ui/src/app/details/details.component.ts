import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ChartConfiguration, ChartOptions } from 'chart.js';
import { Subscription } from 'rxjs';
import { ApiService, Stat } from '../service/api.service';

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

  title = 'ng2-charts-demo';

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
      console.log(this.query)
      this.getDetails();
    })

    this.subscriptions.push(routeSubscription)
  }

  getDetails(): void {
    if (this.query === null) {
      return;
    }
    console.log(this.query)
    let unsafeDetailsSubscription = this.api.getUnsafeDetails(this.query).subscribe( (data: Stat[]) => {
      this.unsafeStats = data
      console.log(data)
      this.mapUnsafeStats()
    })
    this.subscriptions.push(unsafeDetailsSubscription)
  }

  private mapUnsafeStats() {
    let sortedStats = this.unsafeStats.sort((left: Stat, right: Stat) => {
      return left.timestamp - right.timestamp;
    })
    sortedStats = sortedStats.slice(sortedStats.length - 30)
    let values = sortedStats.map(x => x.value);
    let labels = sortedStats.map(x => {
      let t = new Date(1970, 0, 1)
      t.setSeconds(x.timestamp)
      return t
    })
    this.lineChartData.datasets[0].data = values;
    this.lineChartData.labels = labels
    console.log(this.lineChartData)
    this.hasUnsafeStats = true;
  }
}
