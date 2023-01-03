import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core'
import { Observable, throwError } from 'rxjs';
import { catchError, retry } from 'rxjs/operators';

@Injectable({
    providedIn: 'root'
})
export class ApiService {

    private domain: string = 'localhost';
    private port: string = '5321';
    private baseUrl: string; 

    constructor(private http: HttpClient) {
        this.baseUrl = 'http://' + this.domain + ':' + this.port;
    }

    getAllPrompts(): Observable<Prompt[]> {
        return this.http.get<Prompt[]>(this.baseUrl + '/prompts')
    }
}

export interface Prompt {
    query: String,
    isActive: boolean,
    lastReadId: number
  }