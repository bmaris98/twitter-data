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

    togglePrompt(query: string): Observable<Object> {
        return this.http.patch(this.baseUrl + '/prompts/toggle', {
            'query': query,
        });
    }

    addPrompt(query: string): Observable<Object> {
        return this.http.post(this.baseUrl + '/prompts', {
            'query': query,
        });
    }

    getUnsafeDetails(query: string): Observable<Stat[]> {
        console.log(this.baseUrl + '/stats/unsafe/' + query)
        return this.http.get<Stat[]>(this.baseUrl + '/stats/unsafe/' + encodeURIComponent(query))
    }
}

export interface Prompt {
    query: String,
    isActive: boolean,
    lastReadId: number
}

export interface Stat {
    query: String,
    value: number,
    timestamp: number
}