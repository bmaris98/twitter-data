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
        return this.http.get<Stat[]>(this.baseUrl + '/stats/unsafe/' + encodeURIComponent(query))
    }

    getReports(query: string): Observable<Report[]> {
        return this.http.get<Report[]>(this.baseUrl + '/stats/reports/' + encodeURIComponent(query))
    }

    triggerReport(query: string): Observable<Object> {
        return this.http.post(this.baseUrl + '/hadoop/run/' + encodeURIComponent(query), {})
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

export interface Report {
    query: string,
    timestamp: number,
    data: string,
    id: string
}