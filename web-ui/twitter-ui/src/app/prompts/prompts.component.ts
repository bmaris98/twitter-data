import { Component, OnInit } from '@angular/core';
import { Observable, Subscription } from 'rxjs';
import { ApiService, Prompt } from '../service/api.service';

const ELEMENT_DATA: Prompt[] = [
  {query: "@JoeBiden", isActive: false, lastReadId: 0},
  {query: "@elonmusk", isActive: true, lastReadId: 0}
];

@Component({
  selector: 'app-prompts',
  templateUrl: './prompts.component.html',
  styleUrls: ['./prompts.component.scss']
})
export class PromptsComponent implements OnInit {
  displayedColumns: string[] = ['query', 'isActive'];
  dataSource: Prompt[] = [];

  private subscriptions: Subscription[] = []
  public constructor(private api: ApiService) {

  }
  ngOnInit(): void {
    this.api.getAllPrompts().subscribe((data: Prompt[]) => {
      this.dataSource = data;
    })
  }
}