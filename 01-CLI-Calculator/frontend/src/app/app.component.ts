import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  balance: number = 0;
  transactionType: string = 'credit';
  amount: number = 0;
  description: string = '';
  transactions: Transaction[] = [];

  addTransaction() {
    if (this.amount <= 0 || !this.description) {
      return;
    }

    const transaction: Transaction = {
      type: this.transactionType,
      amount: this.amount,
      description: this.description,
      date: new Date()
    };

    this.transactions.unshift(transaction);
    
    if (this.transactionType === 'credit') {
      this.balance += this.amount;
    } else {
      this.balance -= this.amount;
    }

    // Reset form
    this.amount = 0;
    this.description = '';
  }
}

interface Transaction {
  type: string;
  amount: number;
  description: string;
  date: Date;
}