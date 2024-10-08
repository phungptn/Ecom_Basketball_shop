"use client";

import React, { useState, useEffect, useCallback } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import ProductCard from "../../appComoponent/ProductCard";
import { OrderItem } from "../../../lib/productItem";
import {
  convertToOrder,
  FetchOrdersByUserID,
  Order,
  OrderDetailsProps,
  ShippingStatus,
} from "@/lib/order";
import { useOrders, OrdersProvider } from "../../appComoponent/OrdersContext";
import { decryptToken } from "@/lib/decrypt";

export default function OrderDetails() {
  const { orderId } = useParams();
  const { getOrder, updateShippingStatus } = useOrders();
  const [order, setOrder] = useState<Order>();
  //const [order.shippingStatus, setShippingStatus] = useState("delivered");
  const [tokenAvailable, setTokenAvailable] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem("jwt"); // Check immediately if token is available

    if (token) {
      setTokenAvailable(true);
    } else {
      // Set up an interval to keep checking for the token
      const checkForToken = setInterval(() => {
        const token = localStorage.getItem("jwt");
        if (token) {
          setTokenAvailable(true);
          clearInterval(checkForToken);
        }
      }, 100); // Adjust interval as needed

      return () => clearInterval(checkForToken); // Clear interval on component unmount
    }
  }, []);

  useEffect(() => {
    if (!tokenAvailable) return; // Only proceed if token is available

    const token = localStorage.getItem("jwt");
    if (token) {
      const decrypted = decryptToken(token);
      const payload = JSON.parse(atob(decrypted.split(".")[1]));
      const userID = payload.userID;

      const fetchOrder = async () => {
        try {
          const normalizedOrderId = Array.isArray(orderId)
            ? orderId[0]
            : orderId;
          // const order = await getOrder(normalizedOrderId);

          // Fetch orders by user ID
          const orderRequests = await FetchOrdersByUserID(userID);

          // Find the correct order request by order ID
          const matchingOrderRequest = orderRequests
            ? orderRequests.orders.find(
                (orderRequest) => orderRequest.orderID === normalizedOrderId
              )
            : null;

          if (matchingOrderRequest) {
            // Convert the matching order request to Order
            const order = await convertToOrder(matchingOrderRequest);
            setOrder(order);
            console.log("Fetched Order: ", order);
          } else {
            console.error("Order not found with ID: ", normalizedOrderId);
          }
        } catch (error) {
          console.error("Failed to fetch order details: ", error);
          //throw new Error();
        }
      };

      fetchOrder();
    }
  }, [orderId, getOrder]);

  if (order === undefined) {
    return (
      <div className="flex items-center justify-center h-screen">
        <p className="text-xl font-semibold">
          The order ID {orderId} does not exist
        </p>
      </div>
    );
  }

  // Function to render the appropriate stamp based on order.shippingStatus
  const renderShippingStatusStamp = () => {
    switch (order.shippingStatus) {
      case ShippingStatus.Delivered:
        return (
          <p className="successful-delivery-stamp">
            {ShippingStatus.Delivered}
          </p>
        );
      case ShippingStatus.Delivering:
        return (
          <p className="ongoing-delivery-stamp">{ShippingStatus.Delivering}</p>
        );
      case ShippingStatus.Received:
        return (
          <p className="recieved-order-stamp">{ShippingStatus.Received}</p>
        );
      case ShippingStatus.Pending:
        return <p className="pending-order-stamp">{ShippingStatus.Pending}</p>;
      default:
        return null;
    }
  };

  const renderUserButtons = () => {
    if (order.shippingStatus === ShippingStatus.Received) return null;

    return (
      <div className="mt-6 flex justify-between gap-4">
        {/* Cancel button container */}
        {/* {(order.shippingStatus === ShippingStatus.Delivering ||
          order.shippingStatus === ShippingStatus.Pending) && (
          <div className="flex-grow">
            <button
              className="text-error text-xl font-semibold py-4 px-8 rounded-xl
                transition
                transition-duration-300
                transition-property:scale,box-shadow,background-color
                hover:scale-105 hover:drop-shadow-xl hover:bg-accent hover:text-accent-content
                outline-none
                border
                border-error"
              onClick={handleCancelButton}
            >
              Cancel order
            </button>
          </div>
        )} */}

        {/* Confirm button container */}
        {
          /* {order.shippingStatus === ShippingStatus.Delivered &&  }*/
          <div className="ml-auto">
            <button
              className="bg-secondary text-primary-content text-xl font-semibold py-2 px-4
                hover:bg-secondary-content hover:text-secondary
                outline-none
                border-none"
              onClick={handleConfirmButton}
            >
              I have received it!
            </button>
          </div>
        }
      </div>
    );
  };

  const handleConfirmButton = () => {
    updateShippingStatus(order.orderID, ShippingStatus.Received);
  };
  // const handleCancelButton = () => {
  //   updateShippingStatus(order.orderID, ShippingStatus.Canceled);
  // };

  return (
    <div className="p-10">
      <div className="max-w-7xl mx-auto p-6 bg-white shadow-md rounded-lg">
        <h1 className="text-2xl font-semibold mb-6 text-primary-content">
          Order&apos;s ID.
          {<span className="font-semibold text-info pl-2">#{orderId}</span>}
        </h1>
        <div className="space-y-4"></div>
        <div className="flex flex-col gap-6">
          {order.orderItems.map((item, index) => (
            <div
              key={item.productID}
              className="flex flex-row items-center gap-20 px-16 py-6 bg-primary h-full"
            >
              <ProductCard
                key={item.productID}
                product={item}
                isEditable={false}
              />
            </div>
          ))}
        </div>
        <div className="mt-6 text-primary-content">
          <div className="flex justify-between">
            <p className="font-medium text-xl">Coupon applied :</p>
            <p className="font-medium text-lg">{order.coupon}</p>
          </div>

          <div className="flex justify-between mt-2">
            <p className="font-medium text-xl">Delivery method :</p>
            <p className="font-medium text-lg">{order.shippingMethod}</p>
          </div>

          <div className="flex justify-between mt-2">
            <p className="font-medium text-xl">Estimated time of delivery :</p>
            <p className="font-medium text-lg">{order.shipDate}</p>
          </div>

          <div className="flex justify-between mt-2">
            <p className="font-medium text-xl">Payment Method :</p>
            <p className="font-medium text-lg">{order.paymentMethod}</p>
          </div>

          <div className="flex justify-between mt-4">
            <p className="font-bold text-2xl">Total bill :</p>
            <p className="font-bold text-2xl">$ {order.totalBill}</p>
          </div>

          <div className="flex justify-between mt-4">
            <p className="font-medium text-xl">Order status :</p>
            {renderShippingStatusStamp()}
          </div>
        </div>

        {renderUserButtons()}
      </div>
    </div>
  );
}
