import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { PromptsComponent } from './prompts/prompts.component';

const routes: Routes = [
  {path: '', component: PromptsComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {

}
