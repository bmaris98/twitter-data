import { Component, OnInit, OnDestroy } from '@angular/core';
import { MatSlideToggleChange } from '@angular/material/slide-toggle';
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
export class PromptsComponent implements OnInit, OnDestroy {
  displayedColumns: string[] = ['query', 'lastIdRead', 'isActive', 'open'];
  dataSource: Prompt[] = [];

  newPromptValue: string = "";

  private subscriptions: Subscription[] = []
  public constructor(private api: ApiService) {

  }
  ngOnDestroy(): void {
    this.subscriptions.forEach(x => x.unsubscribe())
  }

  ngOnInit(): void {
    this.updatePrompts();
  }

  private updatePrompts() {
    let subscription = this.api.getAllPrompts().subscribe((data: Prompt[]) => {
      this.dataSource = data;
    })
    this.subscriptions.push(subscription);
  }

  onToggle($event: MatSlideToggleChange) {
    console.log($event);
    const query = $event.source.name;
    if (query === null) {
      return;
    }
    let subscription = this.api.togglePrompt(query).subscribe(() => {
      this.updatePrompts();
    })
    this.subscriptions.push(subscription)
  }

  addPromptHandler() {
    console.log(this.newPromptValue)
    let subscription = this.api.addPrompt(this.newPromptValue).subscribe(() => {
      this.updatePrompts();
    })
    this.subscriptions.push(subscription)
    this.newPromptValue = '';
  }
}