import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { DetailsComponent } from './details/details.component';
import { PromptsComponent } from './prompts/prompts.component';

const routes: Routes = [
  {path: 'details/:query', component: DetailsComponent},
  {path: '', component: PromptsComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {

}
